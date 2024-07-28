# Solar-exporter

Данный проэкт осуществляет экспорт метрик инвертора SmartWatt ECO 7.2kW через последовательный порт в prometheus.

## Что нужно для зупкска экспортера?

1) Настроенный prometheus
2) Одноплатный компьютер с архитектурой amd64, i386, arm для установки и запуска экспортера. Идеально подходит Raspberry Pi Zero / Raspberry Pi Zero W2 / NanoPi NEO
3) Кабель USB-RS232 и переходник DB9-8P8C(RJ45). Можно обойтись и без такого кабеля если использовать имеющийся находящийся на плате Serial интерфейс.

### Конфигурация prometheus

```yaml
scrape_configs:
  - job_name: 'smartwatteco'
    scrape_interval: 15s
    static_configs:
      - targets:
        - "/dev/ttyUSB0"
        labels:
          baudrate: "2400"
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_port
      - source_labels: [baudrate]
        target_label: __param_baudrate
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: 192.168.88.45:9560
```
### Пояснения к конфигурации.

- 192.168.88.45 - ip адрес или доменное имя устройства на котором установлкн экспортер
- /dev/ttyUSB0 - имя порта присвоенное адаптеру внутри операционной системы. Можно использовать ID
- baudrate: "2400" - Скорость обмена данными с инвертором. Для модели 7.2к это 2400 бод.

В через одно устройство можно мониторить группу инверторов. Для этого необходимо просто перечислить последовательные порты в targets в конфигурации прометеуса.

### Доступные env

| Назывние | Значение по умолчанию | Описание |
|:---------|:----------------------|:---------|
| ENVPREFIX | solarexporter | Профиск env
| SOLAREXPORTER_CORE_BINDPORT | 9560 | Порт на который повесится экспортер |
| SOLAREXPORTER_CORE_BINDADDR | 0.0.0.0 | Адрес на который вешается экспортер. По умолсчанию на все |
| SOLAREXPORTER_CORE_METRICSPATH | /metrics | Путь в URL по которому будут находится метрики. |
| SOLAREXPORTER_CORE_METRICSPREFIX | solar_invertor | Префикс метрик которые отдаёт экспортер