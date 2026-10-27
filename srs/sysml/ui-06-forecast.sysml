package UI_06_Forecast {
    private import UiKit::*;
    part <forecast> 'Аналитика и прогнозы' : Screen {
        doc /* REQ-FUNC-010, REQ-FUNC-014, REQ-FUNC-015, REQ-FUNC-016, REQ-FUNC-017 */
        part 'Навигация' : NavBar {
            part 'Обзор' : Link;
            part 'Операции' : Link;
            part 'Бюджет' : Link;
            part 'Платежи' : Link;
            part 'Прогноз' : Link;
            part 'Выйти' : Link;
        }
        part 'Сравнение периодов' : Row {
            part 'Период A [01–26.10.2026]' : Input;
            part 'Период Б [01–26.09.2026]' : Input;
            part 'Сравнить' : Button;
        }
        part 'Сравнение: A = 1000; Б = 900 BYN; изменение +100 BYN (+11,11%)' : Label;
        part 'Расходы октября [неполный месяц]' : Chart {
            part 'Продукты 450 (45%); транспорт 150 (15%); жильё 400 (40%)' : Label;
        }
        part 'Прогнозы на 26.10.2026' : Row {
            part 'Остаток 31 октября [850,00 BYN]' : Label;
            part 'Расходы ноября [400,00 BYN]' : Label;
        }
        part 'История и расчёт' : Column {
            part 'Июль–сентябрь: 920 BYN / 92 дня = 10 BYN в день' : Label;
            part 'Остаток: 1000 + 0 − 100 − 10 × 5 = 850 BYN' : Label;
            part 'Расходы ноября: 10 × 30 + 100 = 400 BYN' : Label;
            part 'Связанные регулярные расходы исключены из средней' : Label;
        }
        part 'Состояния' : Column {
            part 'Нет полного месяца — числовые прогнозы не показаны' : Banner;
            part 'Нулевая база сравнения — процент заменён прочерком' : Banner;
            part 'Отрицательный прогноз — возможен дефицит' : Banner;
        }
        part 'Оценка при прежнем уровне трат; операции текущего дня должны быть внесены' : Label;
    }
}
