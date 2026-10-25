"""Generate the requirement model from the versioned requirements register."""
from pathlib import Path
import json
R=Path(__file__).resolve().parents[1]
req=json.loads((R/'srs/data/requirements.json').read_text())
text='''package Balance {
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
'''
groups={'Доступ':'Access','Учёт':'Accounting','Анализ':'Analysis','Планирование':'Planning','Прогноз':'Forecasting','Интерфейсы':'Interfaces','Качество':'Quality','Безопасность':'Security','Ограничения':'Constraints'}
paths={r['id']:groups[r['group']]+'::'+r['id'].replace('-','_') for r in req}
last=None
for r in sorted(req,key=lambda r:list(groups).index(r['group'])):
 if r['group']!=last:
  if last: text+='        }\n'
  text+='        package '+groups[r['group']]+' {\n'
  last=r['group']
 name=r['id'].replace('-','_')
 text+=f"            requirement <'{r['id']}'> {name} : ProductRequirement {{\n                doc /* {r['title']}. {r['statement']} */\n                attribute :>> priority = \"{r['priority']}\";\n            }}\n"
text+='        }\n'
text+='    }\n\n    package Satisfaction {\n'
for r in req:text+=f"        satisfy Requirements::{paths[r['id']]} by system;\n"
text+='    }\n\n    package Scenarios {\n'
ucs=[('Access','Войти в личный бюджет',[1,2]),('RecordOperation','Записать и исправить операцию',[3,4,5,6,7,8]),('InspectBudget','Просмотреть бюджет и аналитику',[9,10]),('SetLimits','Установить лимиты категорий',[11]),('ManagePayments','Запланировать и подтвердить платёж',[12,13,18]),('Forecast','Получить и проверить оба прогноза',[14,15,16,17])]
for name,title,ids in ucs:
 text+=f"        use case def {name} {{\n            doc /* {title}. */\n            subject product : BalanceSystem;\n            actor owner : BudgetOwner;\n            objective {{ doc /* {title}. */ }}\n            }}\n"
 for i in ids:text+=f"        dependency from {name} to Requirements::{paths[f'REQ-FUNC-{i:03d}']};\n"
text+='''    }

    package UserInterface {
        part def Control { attribute label : String; }
        part def InputField :> Control;
        part def Button :> Control;
        part def Screen {
            attribute label : String;
            attribute mockup : String;
        }
'''
controls={'overview':[('addOperation','Button','Добавить операцию')], 'operations':[('amount','InputField','Сумма BYN'),('date','InputField','Дата'),('category','InputField','Категория'),('operationComment','InputField','Комментарий'),('save','Button','Сохранить операцию')], 'budget':[('month','InputField','Месяц'),('limit','InputField','Лимит BYN')], 'scheduled':[('confirm','Button','Подтвердить платёж')], 'forecast':[('history','Control','Период истории'),('monthEnd','Control','Остаток к концу месяца'),('nextMonth','Control','Расходы следующего месяца')]}
screens=[('overview','Обзор бюджета','dashboard'),('operations','Операции','operations'),('budget','Категории и лимиты','budget'),('scheduled','Плановые платежи','scheduled'),('forecast','Аналитика и прогноз','forecast')]
for name,label,file in screens:
 text+=f'        part {name} : Screen {{\n            attribute :>> label = "{label}";\n            attribute :>> mockup = "../mockups/{file}.html";\n'
 for control,kind,lab in controls[name]:text+=f'            part {control} : {kind} {{ attribute :>> label = "{lab}"; }}\n'
 text+='        }\n'
text+='''        item def OpenOperations;
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
'''
(R/'srs/sysml/balance.sysml').write_text(text)
print(f'Model: {len(req)} requirements, {len(req)} satisfy usages, {len(ucs)} use cases')
