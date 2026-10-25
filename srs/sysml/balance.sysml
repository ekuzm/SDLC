package Balance {
    doc /* Баланс. Требования версии 0.2. Связи satisfy обозначают
        проектное назначение, а не результат проверки реализации. */
    private import ScalarValues::*;

    part def BudgetOwner;
    part def BalanceSystem;
    part system : BalanceSystem;

    requirement def ProductRequirement {
        subject product : BalanceSystem;
        attribute priority : String;
    }

    package Requirements {
        package Access {
            requirement <'REQ-FUNC-001'> REQ_FUNC_001 : ProductRequirement {
                doc /* Регистрация. Система должна создавать учётную запись по уникальному адресу электронной почты и паролю. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-FUNC-002'> REQ_FUNC_002 : ProductRequirement {
                doc /* Вход и выход. Система должна предоставлять вход по адресу и паролю и завершать текущий сеанс при выходе. */
                attribute :>> priority = "Must";
            }
        }
        package Accounting {
            requirement <'REQ-FUNC-003'> REQ_FUNC_003 : ProductRequirement {
                doc /* Начало учёта. Система должна хранить для пользователя единый бюджет в BYN, дату начала учёта и начальный остаток. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-FUNC-004'> REQ_FUNC_004 : ProductRequirement {
                doc /* Категории. Система должна позволять создавать, переименовывать и архивировать категории доходов и расходов. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-FUNC-005'> REQ_FUNC_005 : ProductRequirement {
                doc /* Добавление операции. Система должна сохранять доход или расход с суммой, датой, категорией и необязательным комментарием. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-FUNC-006'> REQ_FUNC_006 : ProductRequirement {
                doc /* Изменение операции. Система должна позволять изменять собственную операцию с пересчётом остатка, лимитов, аналитики и прогнозов. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-FUNC-007'> REQ_FUNC_007 : ProductRequirement {
                doc /* Удаление операции. Система должна удалять выбранную собственную операцию только после подтверждения. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-FUNC-008'> REQ_FUNC_008 : ProductRequirement {
                doc /* Поиск и фильтры. Система должна фильтровать операции по включительным границам дат, типу и категории. */
                attribute :>> priority = "Must";
            }
        }
        package Analysis {
            requirement <'REQ-FUNC-009'> REQ_FUNC_009 : ProductRequirement {
                doc /* Обзор бюджета. Система должна показывать текущий остаток, фактические доходы и расходы выбранного месяца и ближайшие плановые платежи. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-FUNC-010'> REQ_FUNC_010 : ProductRequirement {
                doc /* Аналитика расходов. Система должна показывать суммы и доли расходов по категориям и сравнивать два выбранных периода. */
                attribute :>> priority = "Must";
            }
        }
        package Planning {
            requirement <'REQ-FUNC-011'> REQ_FUNC_011 : ProductRequirement {
                doc /* Лимиты категорий. Система должна хранить месячные лимиты расходов по категориям и показывать их использование. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-FUNC-012'> REQ_FUNC_012 : ProductRequirement {
                doc /* Регулярные платежи. Система должна формировать ежемесячные плановые доходы и расходы по заданному пользователем расписанию. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-FUNC-013'> REQ_FUNC_013 : ProductRequirement {
                doc /* Исполнение платежа. Система должна связывать плановое вхождение с одной фактической операцией при подтверждении исполнения. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-FUNC-018'> REQ_FUNC_018 : ProductRequirement {
                doc /* Просроченный плановый платёж. Система должна показывать неисполненное вхождение с прошедшей датой как просроченное и требовать решения пользователя. */
                attribute :>> priority = "Must";
            }
        }
        package Forecasting {
            requirement <'REQ-FUNC-014'> REQ_FUNC_014 : ProductRequirement {
                doc /* Прогноз остатка. Система должна оценивать остаток на конец текущего месяца по правилам F1–F4 раздела 3.2.1. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-FUNC-015'> REQ_FUNC_015 : ProductRequirement {
                doc /* Прогноз расходов. Система должна оценивать расходы следующего календарного месяца по правилу F5. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-FUNC-016'> REQ_FUNC_016 : ProductRequirement {
                doc /* История для прогноза. Система должна выбирать до трёх последних полных календарных месяцев учёта перед текущим месяцем. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-FUNC-017'> REQ_FUNC_017 : ProductRequirement {
                doc /* Объяснение прогноза. Система должна показывать дату расчёта, период истории, формулу и исходные суммы каждого прогноза. */
                attribute :>> priority = "Must";
            }
        }
        package Interfaces {
            requirement <'REQ-INT-001'> REQ_INT_001 : ProductRequirement {
                doc /* Адаптивный интерфейс. Интерфейс должен поддерживать ширину окна от 360 до 1440 пикселей. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-INT-002'> REQ_INT_002 : ProductRequirement {
                doc /* Состояния экранов. Интерфейс должен различать загрузку, отсутствие данных, ошибку ввода и ошибку сохранения. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-INT-003'> REQ_INT_003 : ProductRequirement {
                doc /* Веб-доступ. Система должна предоставлять русскоязычный веб-интерфейс через HTTPS при сетевом развёртывании. */
                attribute :>> priority = "Must";
            }
        }
        package Quality {
            requirement <'REQ-PERF-001'> REQ_PERF_001 : ProductRequirement {
                doc /* Отклик списков. Система должна отвечать на запрос списка операций не более чем за 2 секунды для 95 процентов запросов в контрольном окружении. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-PERF-002'> REQ_PERF_002 : ProductRequirement {
                doc /* Расчёт прогноза. Система должна рассчитывать оба прогноза не более чем за 1 секунду для 95 процентов запросов в контрольном окружении REQ-PERF-001. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-REL-001'> REQ_REL_001 : ProductRequirement {
                doc /* Целостность операций. Система должна сохранять изменение операции и связанных отметок планового платежа атомарно. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-REL-002'> REQ_REL_002 : ProductRequirement {
                doc /* Резервное копирование. Развёрнутая система должна поддерживать ежедневную резервную копию данных и восстановление из неё. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-AVAIL-001'> REQ_AVAIL_001 : ProductRequirement {
                doc /* Контроль работоспособности. Система должна предоставлять проверку готовности принимать запросы. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-OBS-001'> REQ_OBS_001 : ProductRequirement {
                doc /* Диагностика ошибок. Система должна регистрировать технические ошибки с временем и идентификатором запроса. */
                attribute :>> priority = "Must";
            }
        }
        package Security {
            requirement <'REQ-SEC-001'> REQ_SEC_001 : ProductRequirement {
                doc /* Защита паролей. Система должна хранить только стойкие солёные хеши паролей и ограничивать частоту попыток входа. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-SEC-002'> REQ_SEC_002 : ProductRequirement {
                doc /* Защита сеанса. Система должна защищать сеанс пользователя и изменения данных от межсайтовой подделки запросов. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-SEC-003'> REQ_SEC_003 : ProductRequirement {
                doc /* Изоляция данных. Система должна проверять владельца при каждом чтении и изменении бюджетных данных. */
                attribute :>> priority = "Must";
            }
        }
        package Constraints {
            requirement <'REQ-COMP-001'> REQ_COMP_001 : ProductRequirement {
                doc /* Прозрачность назначения. Система должна обозначать прогноз как информационную оценку, основанную на введённых пользователем данных. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-INST-001'> REQ_INST_001 : ProductRequirement {
                doc /* Воспроизводимый запуск. Проект должен содержать инструкцию установки и запуска серверной части и хранилища. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-BUILD-001'> REQ_BUILD_001 : ProductRequirement {
                doc /* Воспроизводимая сборка. Серверная часть должна собираться средствами Go с зафиксированными версиями зависимостей. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-MAINT-001'> REQ_MAINT_001 : ProductRequirement {
                doc /* Проверка финансовой логики. Проект должен содержать автоматические проверки остатка, лимитов, регулярных платежей и обеих формул прогноза. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-PORT-001'> REQ_PORT_001 : ProductRequirement {
                doc /* Технологические ограничения. Серверная логика должна быть реализована на Go, а пользовательский доступ — через браузер. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-COST-001'> REQ_COST_001 : ProductRequirement {
                doc /* Доступность демонстрации. Локальная демонстрация обязательных функций должна выполняться без платных API и подписок. */
                attribute :>> priority = "Must";
            }
            requirement <'REQ-CM-001'> REQ_CM_001 : ProductRequirement {
                doc /* Управление требованиями. Каждое изменение требования должно сохранять его идентификатор и фиксироваться в истории документа и Git. */
                attribute :>> priority = "Must";
            }
        }
    }

    package Satisfaction {
        satisfy Requirements::Access::REQ_FUNC_001 by system;
        satisfy Requirements::Access::REQ_FUNC_002 by system;
        satisfy Requirements::Accounting::REQ_FUNC_003 by system;
        satisfy Requirements::Accounting::REQ_FUNC_004 by system;
        satisfy Requirements::Accounting::REQ_FUNC_005 by system;
        satisfy Requirements::Accounting::REQ_FUNC_006 by system;
        satisfy Requirements::Accounting::REQ_FUNC_007 by system;
        satisfy Requirements::Accounting::REQ_FUNC_008 by system;
        satisfy Requirements::Analysis::REQ_FUNC_009 by system;
        satisfy Requirements::Analysis::REQ_FUNC_010 by system;
        satisfy Requirements::Planning::REQ_FUNC_011 by system;
        satisfy Requirements::Planning::REQ_FUNC_012 by system;
        satisfy Requirements::Planning::REQ_FUNC_013 by system;
        satisfy Requirements::Forecasting::REQ_FUNC_014 by system;
        satisfy Requirements::Forecasting::REQ_FUNC_015 by system;
        satisfy Requirements::Forecasting::REQ_FUNC_016 by system;
        satisfy Requirements::Forecasting::REQ_FUNC_017 by system;
        satisfy Requirements::Planning::REQ_FUNC_018 by system;
        satisfy Requirements::Interfaces::REQ_INT_001 by system;
        satisfy Requirements::Interfaces::REQ_INT_002 by system;
        satisfy Requirements::Interfaces::REQ_INT_003 by system;
        satisfy Requirements::Quality::REQ_PERF_001 by system;
        satisfy Requirements::Quality::REQ_PERF_002 by system;
        satisfy Requirements::Security::REQ_SEC_001 by system;
        satisfy Requirements::Security::REQ_SEC_002 by system;
        satisfy Requirements::Security::REQ_SEC_003 by system;
        satisfy Requirements::Quality::REQ_REL_001 by system;
        satisfy Requirements::Quality::REQ_REL_002 by system;
        satisfy Requirements::Quality::REQ_AVAIL_001 by system;
        satisfy Requirements::Quality::REQ_OBS_001 by system;
        satisfy Requirements::Constraints::REQ_COMP_001 by system;
        satisfy Requirements::Constraints::REQ_INST_001 by system;
        satisfy Requirements::Constraints::REQ_BUILD_001 by system;
        satisfy Requirements::Constraints::REQ_MAINT_001 by system;
        satisfy Requirements::Constraints::REQ_PORT_001 by system;
        satisfy Requirements::Constraints::REQ_COST_001 by system;
        satisfy Requirements::Constraints::REQ_CM_001 by system;
    }

    package Scenarios {
        use case def Access {
            doc /* Войти в личный бюджет. */
            subject product : BalanceSystem;
            actor owner : BudgetOwner;
            objective { doc /* Войти в личный бюджет. */ }
            }
        dependency from Access to Requirements::Access::REQ_FUNC_001;
        dependency from Access to Requirements::Access::REQ_FUNC_002;
        use case def RecordOperation {
            doc /* Записать и исправить операцию. */
            subject product : BalanceSystem;
            actor owner : BudgetOwner;
            objective { doc /* Записать и исправить операцию. */ }
            }
        dependency from RecordOperation to Requirements::Accounting::REQ_FUNC_003;
        dependency from RecordOperation to Requirements::Accounting::REQ_FUNC_004;
        dependency from RecordOperation to Requirements::Accounting::REQ_FUNC_005;
        dependency from RecordOperation to Requirements::Accounting::REQ_FUNC_006;
        dependency from RecordOperation to Requirements::Accounting::REQ_FUNC_007;
        dependency from RecordOperation to Requirements::Accounting::REQ_FUNC_008;
        use case def InspectBudget {
            doc /* Просмотреть бюджет и аналитику. */
            subject product : BalanceSystem;
            actor owner : BudgetOwner;
            objective { doc /* Просмотреть бюджет и аналитику. */ }
            }
        dependency from InspectBudget to Requirements::Analysis::REQ_FUNC_009;
        dependency from InspectBudget to Requirements::Analysis::REQ_FUNC_010;
        use case def SetLimits {
            doc /* Установить лимиты категорий. */
            subject product : BalanceSystem;
            actor owner : BudgetOwner;
            objective { doc /* Установить лимиты категорий. */ }
            }
        dependency from SetLimits to Requirements::Planning::REQ_FUNC_011;
        use case def ManagePayments {
            doc /* Запланировать и подтвердить платёж. */
            subject product : BalanceSystem;
            actor owner : BudgetOwner;
            objective { doc /* Запланировать и подтвердить платёж. */ }
            }
        dependency from ManagePayments to Requirements::Planning::REQ_FUNC_012;
        dependency from ManagePayments to Requirements::Planning::REQ_FUNC_013;
        dependency from ManagePayments to Requirements::Planning::REQ_FUNC_018;
        use case def Forecast {
            doc /* Получить и проверить оба прогноза. */
            subject product : BalanceSystem;
            actor owner : BudgetOwner;
            objective { doc /* Получить и проверить оба прогноза. */ }
            }
        dependency from Forecast to Requirements::Forecasting::REQ_FUNC_014;
        dependency from Forecast to Requirements::Forecasting::REQ_FUNC_015;
        dependency from Forecast to Requirements::Forecasting::REQ_FUNC_016;
        dependency from Forecast to Requirements::Forecasting::REQ_FUNC_017;
    }

    package UserInterface {
        part def Control { attribute label : String; }
        part def InputField :> Control;
        part def Button :> Control;
        part def Screen {
            attribute label : String;
            attribute mockup : String;
        }
        part overview : Screen {
            attribute :>> label = "Обзор бюджета";
            attribute :>> mockup = "../mockups/dashboard.html";
            part addOperation : Button { attribute :>> label = "Добавить операцию"; }
        }
        part operations : Screen {
            attribute :>> label = "Операции";
            attribute :>> mockup = "../mockups/operations.html";
            part amount : InputField { attribute :>> label = "Сумма BYN"; }
            part date : InputField { attribute :>> label = "Дата"; }
            part category : InputField { attribute :>> label = "Категория"; }
            part operationComment : InputField { attribute :>> label = "Комментарий"; }
            part save : Button { attribute :>> label = "Сохранить операцию"; }
        }
        part budget : Screen {
            attribute :>> label = "Категории и лимиты";
            attribute :>> mockup = "../mockups/budget.html";
            part month : InputField { attribute :>> label = "Месяц"; }
            part limit : InputField { attribute :>> label = "Лимит BYN"; }
        }
        part scheduled : Screen {
            attribute :>> label = "Плановые платежи";
            attribute :>> mockup = "../mockups/scheduled.html";
            part confirm : Button { attribute :>> label = "Подтвердить платёж"; }
        }
        part forecast : Screen {
            attribute :>> label = "Аналитика и прогноз";
            attribute :>> mockup = "../mockups/forecast.html";
            part history : Control { attribute :>> label = "Период истории"; }
            part monthEnd : Control { attribute :>> label = "Остаток к концу месяца"; }
            part nextMonth : Control { attribute :>> label = "Расходы следующего месяца"; }
        }
        item def OpenOperations;
        item def OpenBudget;
        item def OpenScheduled;
        item def OpenForecast;
        item def OpenOverview;
        item def SubmitOperation;
        item def CorrectInput;
        item def SubmitValidOperation;
        state def Navigation {
            entry; then dashboard;
            state dashboard;
            state operationList;
            state budgetLimits;
            state scheduledPayments;
            state analysisScreen;
            transition first dashboard accept OpenOperations then operationList;
            transition first dashboard accept OpenBudget then budgetLimits;
            transition first dashboard accept OpenScheduled then scheduledPayments;
            transition first dashboard accept OpenForecast then analysisScreen;
            transition first operationList accept OpenOverview then dashboard;
            transition first budgetLimits accept OpenOverview then dashboard;
            transition first scheduledPayments accept OpenOverview then dashboard;
            transition first analysisScreen accept OpenOverview then dashboard;
        }
        state def OperationForm {
            entry; then editing;
            state editing;
            state invalidInput;
            state saved;
            transition first editing accept SubmitOperation then invalidInput;
            transition first invalidInput accept CorrectInput then editing;
            transition first editing accept SubmitValidOperation then saved;
        }
        dependency from overview to Scenarios::InspectBudget;
        dependency from operations to Scenarios::RecordOperation;
        dependency from budget to Scenarios::SetLimits;
        dependency from scheduled to Scenarios::ManagePayments;
        dependency from forecast to Scenarios::Forecast;
    }

    package Domain {
        item def Category { attribute name : String; attribute kind : String; attribute archived : Boolean; }
        item def Operation {
            attribute amountKopecks : Integer;
            attribute date : String;
            ref item category : Category;
        }
        item def PlannedOccurrence {
            attribute dueDate : String;
            attribute status : String;
            ref item operation : Operation[0..1];
        }
        item def RecurringRule {
            attribute dayOfMonth : Integer;
            item occurrences : PlannedOccurrence[0..*];
        }
        item def CategoryLimit {
            attribute month : String;
            attribute amountKopecks : Integer;
            ref item category : Category;
        }
        item def Budget {
            attribute startDate : String;
            attribute initialKopecks : Integer;
            item categories : Category[0..*];
            item operations : Operation[0..*];
            item limits : CategoryLimit[0..*];
            item rules : RecurringRule[0..*];
        }
        item def User {
            attribute email : String;
            item budget : Budget;
        }
    }
}
