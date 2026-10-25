"""Generate the requirement model from the versioned requirements register."""
from pathlib import Path
import json
R=Path(__file__).resolve().parents[1]
req=json.loads((R/'srs/data/requirements.json').read_text())
text='''package Balance {
    doc /* Баланс. Требования версии 0.3. Связи satisfy обозначают
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
text+='    }\n}\n'
(R/'srs/sysml/balance.sysml').write_text(text)
print(f'Model: {len(req)} requirements and {len(req)} satisfy usages')
