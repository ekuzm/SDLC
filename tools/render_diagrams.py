"""Readable SVG projections of the SysML source; PNG conversion uses sharp."""
from pathlib import Path
import re,html,textwrap
R=Path(__file__).resolve().parents[1];source=(R/'srs/sysml/balance.sysml').read_text();D=R/'srs/sysml/views';D.mkdir(exist_ok=True)
def esc(x):return html.escape(str(x))
def canvas(w,h,title,sub):return [f'<svg xmlns="http://www.w3.org/2000/svg" width="{w}" height="{h}" viewBox="0 0 {w} {h}"><defs><marker id="arrow" markerWidth="10" markerHeight="10" refX="9" refY="5" orient="auto"><path d="M0 0 L10 5 L0 10 Z" fill="#78919f"/></marker></defs><rect width="100%" height="100%" fill="white"/>',f'<text x="36" y="44" font-family="Arial" font-size="25" font-weight="bold" fill="#17394b">{esc(title)}</text>',f'<text x="36" y="72" font-family="Arial" font-size="14" fill="#617c8b">{esc(sub)}</text>']
def label(a,x,y,text,size=17,bold=False,fill='#203b4c',anchor='start'):
 a.append(f'<text x="{x}" y="{y}" font-family="Arial" font-size="{size}" font-weight="{"bold" if bold else "normal"}" fill="{fill}" text-anchor="{anchor}">{esc(text)}</text>')
def box(a,x,y,w,h,typ,title,extra='',tint=False):
 a.append(f'<rect x="{x}" y="{y}" width="{w}" height="{h}" rx="8" fill="{"#eaf5f0" if tint else "#f5f8fa"}" stroke="#afc5ce"/>')
 label(a,x+14,y+23,typ,12,fill='#64808d');label(a,x+14,y+48,title,17,True)
 for j,line in enumerate(textwrap.wrap(extra,width=max(18,int(w/9)))):label(a,x+14,y+70+j*20,line,14)
def edge(a,x,y,xx,yy,text=''):
 a.append(f'<path d="M{x} {y} L{xx} {yy}" fill="none" stroke="#78919f" stroke-width="1.5" marker-end="url(#arrow)"/>')
 if text:label(a,(x+xx)/2+5,(y+yy)/2-8,text,12)
def save(a,name):a.append('</svg>');(D/f'{name}.svg').write_text('\n'.join(a))
groups={'Access':'Доступ','Accounting':'Учёт','Analysis':'Анализ','Planning':'Планирование','Forecasting':'Прогноз','Interfaces':'Интерфейсы','Quality':'Качество','Security':'Безопасность','Constraints':'Ограничения'}
entries={}
for g in groups:
 block=re.search(r'package '+g+r' \{(.*?)\n        \}',source,re.S).group(1)
 entries[g]=re.findall(r"requirement <'([^']+)'>.*?doc /\* (.*?)\. ",block,re.S)
for filename,keys,title in [('functional-requirements',list(groups)[:5],'Функциональные требования'),('quality-requirements',list(groups)[5:],'Нефункциональные требования')]:
 count=sum(len(entries[g]) for g in keys);h=115+count*65+len(keys)*20
 a=canvas(1120,h,title,'Проекция пакетов Balance::Requirements · текстовая модель SysML v2 · SRS 0.2')
 y=110
 for g in keys:
  mid=y+len(entries[g])*65/2
  box(a,30,mid-34,235,65,'package',groups[g])
  for ident,name in entries[g]:
   edge(a,265,mid,335,y+26)
   a.append(f'<rect x="335" y="{y}" width="745" height="53" rx="6" fill="{"#eaf5f0" if ident=="REQ-FUNC-018" else "#f7f9fa"}" stroke="#b5c8d0"/>')
   label(a,350,y+20,ident,13,True,fill='#1c776c');label(a,350,y+41,name,17)
   y+=65
  y+=20
 save(a,filename)
a=canvas(1120,740,'Участник и варианты использования','Balance::Scenarios · subject BalanceSystem · actor BudgetOwner')
box(a,36,300,220,115,'part def','BudgetOwner','Владелец бюджета',True)
usecases=re.findall(r'use case def (\w+) \{\s*doc /\* (.*?)\. \*/',source)
for j,(name,title) in enumerate(usecases):
 y=108+j*95;edge(a,256,356,360,y+38);box(a,360,y,720,78,'use case def '+name,title)
save(a,'use-cases')
a=canvas(1120,780,'Экраны и переходы','Balance::UserInterface::Navigation · события открытия экранов')
box(a,390,115,340,100,'state','dashboard','Обзор бюджета',True)
positions=[(35,350,'operationList','Операции'),(580,350,'budgetLimits','Категории и лимиты'),(35,570,'scheduledPayments','Плановые платежи'),(580,570,'analysisScreen','Аналитика и прогноз')]
for x,y,name,title in positions:
 assert 'state '+name+';' in source
 box(a,x,y,495,110,'state '+name,title)
 if y==350: edge(a,550,215,x+245,y,'Открыть экран')
 else:
  track=16 if x==35 else 1100
  a.append(f'<path d="M{410 if x==35 else 710} 215 L{track} 270 L{track} {y+55} L{x if x==35 else x+495} {y+55}" fill="none" stroke="#78919f" stroke-width="1.5" marker-end="url(#arrow)"/>')
 label(a,x+18,y+91,'OpenOverview → обзор',14,fill='#217c70')
label(a,36,753,'Форма операции: editing → invalidInput → editing → saved',17)
save(a,'navigation')
a=canvas(1120,1010,'Модель предметной области','Balance::Domain · логические данные и связи · без привязки к СУБД')
items=[(35,110,300,110,'User','email; budget[1]'),(420,110,650,110,'Budget','startDate; initialKopecks; единый бюджет в BYN'),(35,350,300,140,'Category','name; тип доход или расход'),(420,350,300,140,'Operation','amountKopecks; date; category'),(770,350,300,140,'CategoryLimit','month; amountKopecks; category'),(35,690,300,150,'RecurringRule','dayOfMonth; occurrences[0..*]'),(420,690,650,150,'PlannedOccurrence','dueDate; status; operation[0..1]')]
for x,y,w,h,name,extra in items:
 assert 'item def '+name in source
 box(a,x,y,w,h,'item def',name,extra,True if name=='Budget' else False)
edge(a,335,162,420,162,'1');edge(a,480,220,185,350,'0..*');edge(a,620,220,570,350,'0..*');edge(a,900,220,920,350,'0..*');a.append('<path d="M420 180 L370 250 L370 650 L185 690" fill="none" stroke="#78919f" marker-end="url(#arrow)"/>');edge(a,335,760,420,760,'0..*');edge(a,720,690,600,490,'0..1');edge(a,420,415,335,415,'1')
label(a,36,911,'Инварианты',19,True);label(a,36,945,'Операция и плановое вхождение связываются один к одному.',17);label(a,36,975,'Начальный остаток и планы не являются фактическими доходами.',17)
save(a,'domain')
print('5 SVG projections built from model identifiers')
