package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	defaultBroker   = "47.92.253.145:1883"
	defaultUsername = "0893d019-c63b-48b3-b40"
	defaultPassword = "c49bd8e"
	defaultTopic    = "devices/telemetry"
	defaultClientID = "mqtt_112f3ddf-587"
	defaultSN       = "00501521042000019454"
)

type config struct {
	broker   string
	username string
	password string
	clientID string
	sn       string
	topic    string
	interval time.Duration
	limits   alarmLimits
}

type telemetry struct {
	Alarm       int     `json:"Alarm,omitempty"`
	LocalTime   string  `json:"sys_local_time"`
	SN          string  `json:"SN"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Illuminance int     `json:"illuminance"`
}

type alarmLimits struct {
	temperatureMin float64
	temperatureMax float64
	humidityMin    float64
	humidityMax    float64
	illuminanceMax float64
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	client, err := connect(cfg)
	if err != nil {
		log.Fatalf("连接 MQTT Broker 失败: %v", err)
	}
	defer client.Disconnect(250)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Printf("USR-PT170 虚拟传感器已启动，Broker=%s，Topic=%s，上报间隔=%s", cfg.broker, cfg.topic, cfg.interval)
	if err := run(ctx, client, cfg, rand.New(rand.NewSource(time.Now().UnixNano()))); err != nil {
		log.Fatal(err)
	}
	log.Println("USR-PT170 虚拟传感器已停止")
}

func loadConfig() (config, error) {
	cfg := config{}
	flag.StringVar(&cfg.broker, "broker", envOrDefault("USR_PT170_MQTT_BROKER", defaultBroker), "MQTT Broker 地址")
	flag.StringVar(&cfg.username, "username", envOrDefault("USR_PT170_MQTT_USERNAME", defaultUsername), "MQTT 用户名")
	flag.StringVar(&cfg.password, "password", envOrDefault("USR_PT170_MQTT_PASSWORD", defaultPassword), "MQTT 密码")
	flag.StringVar(&cfg.clientID, "client-id", envOrDefault("USR_PT170_MQTT_CLIENT_ID", defaultClientID), "MQTT ClientID")
	flag.StringVar(&cfg.sn, "sn", envOrDefault("USR_PT170_SN", defaultSN), "设备 SN")
	flag.StringVar(&cfg.topic, "topic", envOrDefault("USR_PT170_MQTT_TOPIC", defaultTopic), "遥测上报主题")
	flag.DurationVar(&cfg.interval, "interval", envDurationOrDefault("USR_PT170_INTERVAL", 10*time.Second), "遥测上报间隔")
	flag.Float64Var(&cfg.limits.temperatureMin, "temperature-min", 0, "温度报警下限")
	flag.Float64Var(&cfg.limits.temperatureMax, "temperature-max", 35, "温度报警上限")
	flag.Float64Var(&cfg.limits.humidityMin, "humidity-min", 20, "湿度报警下限")
	flag.Float64Var(&cfg.limits.humidityMax, "humidity-max", 85, "湿度报警上限")
	flag.Float64Var(&cfg.limits.illuminanceMax, "illuminance-max", 50000, "光照度报警上限")
	flag.Parse()

	cfg.broker = normalizeBroker(cfg.broker)
	if cfg.username == "" {
		return config{}, errors.New("MQTT 用户名不能为空")
	}
	if cfg.clientID == "" {
		return config{}, errors.New("MQTT ClientID 不能为空")
	}
	if cfg.sn == "" {
		return config{}, errors.New("设备 SN 不能为空")
	}
	if cfg.topic == "" {
		return config{}, errors.New("MQTT Topic 不能为空")
	}
	if cfg.interval <= 0 {
		return config{}, errors.New("上报间隔必须大于 0")
	}
	return cfg, nil
}

func connect(cfg config) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions().
		AddBroker(cfg.broker).
		SetUsername(cfg.username).
		SetPassword(cfg.password).
		SetClientID(cfg.clientID).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second).
		SetMaxReconnectInterval(30 * time.Second).
		SetConnectTimeout(10 * time.Second).
		SetWriteTimeout(10 * time.Second).
		SetOrderMatters(false)

	opts.SetOnConnectHandler(func(mqtt.Client) {
		log.Println("MQTT 连接成功")
	})
	opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		log.Printf("MQTT 连接断开，将自动重连: %v", err)
	})

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if !token.WaitTimeout(15 * time.Second) {
		return nil, errors.New("连接 MQTT Broker 超时")
	}
	if err := token.Error(); err != nil {
		return nil, err
	}
	return client, nil
}

func run(ctx context.Context, client mqtt.Client, cfg config, random *rand.Rand) error {
	startedAt := time.Now()
	if err := publish(client, cfg.topic, sample(startedAt, startedAt, cfg.sn, cfg.limits, random)); err != nil {
		return err
	}

	ticker := time.NewTicker(cfg.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case now := <-ticker.C:
			if err := publish(client, cfg.topic, sample(now, startedAt, cfg.sn, cfg.limits, random)); err != nil {
				log.Printf("遥测上报失败: %v", err)
			}
		}
	}
}

func sample(now, startedAt time.Time, sn string, limits alarmLimits, random *rand.Rand) telemetry {
	seconds := now.Sub(startedAt).Seconds()
	temperature := 24 + 3*math.Sin(seconds/90) + jitter(random, 0.3)
	humidity := 55 - 8*math.Sin(seconds/110) + jitter(random, 1.0)
	illuminance := 10000 + 5000*math.Sin(seconds/70) + jitter(random, 500)
	temperature = round(clamp(temperature, -40, 85), 1)
	humidity = round(clamp(humidity, 0, 100), 1)
	illuminance = math.Round(clamp(illuminance, 0, 100000))

	data := telemetry{
		LocalTime:   now.Format("2006-01-02,15:04:05"),
		SN:          sn,
		Temperature: temperature,
		Humidity:    humidity,
		Illuminance: int(illuminance),
	}
	if isAlarm(temperature, humidity, illuminance, limits) {
		data.Alarm = 1
	}
	return data
}

func isAlarm(temperature, humidity, illuminance float64, limits alarmLimits) bool {
	return temperature < limits.temperatureMin || temperature > limits.temperatureMax ||
		humidity < limits.humidityMin || humidity > limits.humidityMax ||
		illuminance > limits.illuminanceMax
}

func publish(client mqtt.Client, topic string, data telemetry) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("序列化遥测数据失败: %w", err)
	}
	token := client.Publish(topic, 1, false, payload)
	if !token.WaitTimeout(10 * time.Second) {
		return errors.New("发布 MQTT 消息超时")
	}
	if err := token.Error(); err != nil {
		return fmt.Errorf("发布 MQTT 消息失败: %w", err)
	}
	log.Printf("遥测上报成功: %s", payload)
	return nil
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func envDurationOrDefault(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Printf("环境变量 %s=%q 不是有效时长，使用默认值 %s", key, value, fallback)
		return fallback
	}
	return duration
}

func normalizeBroker(broker string) string {
	broker = strings.TrimSpace(broker)
	if strings.Contains(broker, "://") {
		return broker
	}
	return "tcp://" + broker
}

func jitter(random *rand.Rand, amplitude float64) float64 {
	return (random.Float64()*2 - 1) * amplitude
}

func clamp(value, min, max float64) float64 {
	return math.Max(min, math.Min(max, value))
}

func round(value float64, precision int) float64 {
	factor := math.Pow10(precision)
	return math.Round(value*factor) / factor
}
