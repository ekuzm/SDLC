# Графические материалы SRS на SysML v2

Модели для Приложения А [спецификации](../srs.md). PNG находятся в [images](../images/).

| Исходник | Назначение | Изображение |
|---|---|---|
| [balance.sysml](balance.sysml) | 37 требований, иерархия групп и 37 связей satisfy | — |
| [context.sysml](context.sysml) | Контекст и потоки данных | [context.png](../images/context.png) |
| [use-cases.sysml](use-cases.sysml) | Гость, владелец бюджета, восемь сценариев; связи с требованиями | [use-cases.png](../images/use-cases.png) |
| [ui-kit.sysml](ui-kit.sysml) | Общие элементы интерфейса | — |
| [ui-01-auth.sysml](ui-01-auth.sysml) | Вход и регистрация | [ui-01-auth.png](../images/ui-01-auth.png) |
| [ui-02-overview.sysml](ui-02-overview.sysml) | Обзор и начало учёта | [ui-02-overview.png](../images/ui-02-overview.png) |
| [ui-03-operations.sysml](ui-03-operations.sysml) | Операции, фильтры, формы и ошибки | [ui-03-operations.png](../images/ui-03-operations.png) |
| [ui-04-budget.sysml](ui-04-budget.sysml) | Категории и лимиты | [ui-04-budget.png](../images/ui-04-budget.png) |
| [ui-05-payments.sysml](ui-05-payments.sysml) | Регулярные платежи и исполнение | [ui-05-payments.png](../images/ui-05-payments.png) |
| [ui-06-forecast.sysml](ui-06-forecast.sysml) | Аналитика и оба прогноза | [ui-06-forecast.png](../images/ui-06-forecast.png) |

## Построение изображений

Из корня репозитория, при установленном и запущенном Docker:

```sh
sh srs/sysml/tools/render-all.sh
python3 tools/check_artifacts.py
```

Окружение: SysML v2 Pilot Implementation 0.62.0, Java 21, Graphviz. При первом запуске загружаются образ и зависимости. Все модели загружаются одной единицей в Jupyter-ядро SysML для разрешения ссылок между файлами. Ошибка модели или отсутствие результата прерывают экспорт.

`%viz --view=interconnection` строит контекст и экраны; `%viz --view=case` — варианты использования. Для экспорта используется PlantUML, выданный самим ядром. Невидимые связи задают горизонтальное расположение элементов Row и вертикальный порядок блоков экрана. Шрифт DejaVu Sans обеспечивает одинаковые метрики при построении и растеризации. Источник подхода — [PR-образец](https://github.com/kiryshabutor/SDLC/pull/1).

[Манифест](tools/render-manifest.json) содержит хеши исходников и PNG последнего успешного построения. Готовые изображения сохранены в Git и доступны без Docker. Связи satisfy обозначают назначение проектируемой системы; испытания приложения выполняются после его реализации.
