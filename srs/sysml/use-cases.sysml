package BalanceActors {
    part def Guest;
    part def Owner;
}
package BalanceUseCases {
    private import BalanceActors::*;
    use case def <Register> 'Зарегистрироваться' {
        doc /* REQ-FUNC-001 */
        subject product : Balance::BalanceSystem;
        actor 'Гость' : Guest;
    }
    use case def <Login> 'Войти и выйти' {
        doc /* REQ-FUNC-002 */
        subject product : Balance::BalanceSystem;
        actor 'Гость' : Guest;
        actor 'Владелец бюджета' : Owner;
    }
    use case def <Setup> 'Задать начало учёта' {
        doc /* REQ-FUNC-003 */
        subject product : Balance::BalanceSystem;
        actor 'Владелец бюджета' : Owner;
    }
    use case def <Operations> 'Вести операции' {
        doc /* REQ-FUNC-005, REQ-FUNC-006, REQ-FUNC-007, REQ-FUNC-008 */
        subject product : Balance::BalanceSystem;
        actor 'Владелец бюджета' : Owner;
    }
    use case def <Categories> 'Управлять категориями и лимитами' {
        doc /* REQ-FUNC-004, REQ-FUNC-011 */
        subject product : Balance::BalanceSystem;
        actor 'Владелец бюджета' : Owner;
    }
    use case def <Overview> 'Просмотреть бюджет и аналитику' {
        doc /* REQ-FUNC-009, REQ-FUNC-010 */
        subject product : Balance::BalanceSystem;
        actor 'Владелец бюджета' : Owner;
    }
    use case def <Plans> 'Планировать и подтверждать платежи' {
        doc /* REQ-FUNC-012, REQ-FUNC-013, REQ-FUNC-018 */
        subject product : Balance::BalanceSystem;
        actor 'Владелец бюджета' : Owner;
    }
    use case def <Forecast> 'Получить и проверить прогнозы' {
        doc /* REQ-FUNC-014, REQ-FUNC-015, REQ-FUNC-016, REQ-FUNC-017 */
        subject product : Balance::BalanceSystem;
        actor 'Владелец бюджета' : Owner;
    }
}
package BalanceTraceability {
    dependency from BalanceUseCases::'Зарегистрироваться' to Balance::Requirements::Access::REQ_FUNC_001;
    dependency from BalanceUseCases::'Войти и выйти' to Balance::Requirements::Access::REQ_FUNC_002;
    dependency from BalanceUseCases::'Задать начало учёта' to Balance::Requirements::Accounting::REQ_FUNC_003;
    dependency from BalanceUseCases::'Вести операции' to Balance::Requirements::Accounting::REQ_FUNC_005;
    dependency from BalanceUseCases::'Вести операции' to Balance::Requirements::Accounting::REQ_FUNC_006;
    dependency from BalanceUseCases::'Вести операции' to Balance::Requirements::Accounting::REQ_FUNC_007;
    dependency from BalanceUseCases::'Вести операции' to Balance::Requirements::Accounting::REQ_FUNC_008;
    dependency from BalanceUseCases::'Управлять категориями и лимитами' to Balance::Requirements::Accounting::REQ_FUNC_004;
    dependency from BalanceUseCases::'Управлять категориями и лимитами' to Balance::Requirements::Planning::REQ_FUNC_011;
    dependency from BalanceUseCases::'Просмотреть бюджет и аналитику' to Balance::Requirements::Analysis::REQ_FUNC_009;
    dependency from BalanceUseCases::'Просмотреть бюджет и аналитику' to Balance::Requirements::Analysis::REQ_FUNC_010;
    dependency from BalanceUseCases::'Планировать и подтверждать платежи' to Balance::Requirements::Planning::REQ_FUNC_012;
    dependency from BalanceUseCases::'Планировать и подтверждать платежи' to Balance::Requirements::Planning::REQ_FUNC_013;
    dependency from BalanceUseCases::'Планировать и подтверждать платежи' to Balance::Requirements::Planning::REQ_FUNC_018;
    dependency from BalanceUseCases::'Получить и проверить прогнозы' to Balance::Requirements::Forecasting::REQ_FUNC_014;
    dependency from BalanceUseCases::'Получить и проверить прогнозы' to Balance::Requirements::Forecasting::REQ_FUNC_015;
    dependency from BalanceUseCases::'Получить и проверить прогнозы' to Balance::Requirements::Forecasting::REQ_FUNC_016;
    dependency from BalanceUseCases::'Получить и проверить прогнозы' to Balance::Requirements::Forecasting::REQ_FUNC_017;
}
