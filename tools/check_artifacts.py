"""Validate documentation consistency and specification arithmetic, not application code."""
from pathlib import Path
from decimal import Decimal, ROUND_HALF_UP
from datetime import date
import json,re,calendar
R=Path(__file__).resolve().parents[1];S=R/'srs'
req=json.loads((S/'data/requirements.json').read_text());srs=(S/'srs.md').read_text();model=(S/'sysml/balance.sysml').read_text();scenarios=(S/'verification/scenarios.md').read_text()
ids=[r['id'] for r in req];assert len(ids)==len(set(ids))==37
headings=re.findall(r'^#### (REQ-\w+-\d{3}) ',srs,re.M);assert set(ids)==set(headings) and len(headings)==37
model_ids=re.findall(r"requirement <'(REQ-[^']+)'>",model);assert set(ids)==set(model_ids) and len(model_ids)==37
assert len(re.findall(r'\bsatisfy Requirements::',model))==37
for r in req:
 assert '## V-'+r['id'][4:]+' ' in scenarios
 assert f"| {r['id']} |" in srs
 assert r['statement'] in srs and r['statement'] in model
 assert all(c in srs and c in scenarios for c in r['criteria'])
for file in [S/'srs.md',S/'review/interview.md',S/'review/review.md',S/'sysml/README.md']:
 for link in re.findall(r'\]\(([^)]+)\)',file.read_text()):
  if not link.startswith(('https://','http://','#')):assert (file.parent/link.split('#')[0]).exists(),(file,link)
for fig in ['functional-requirements','quality-requirements']:
 svg=(S/f'sysml/views/{fig}.svg').read_text()
 expected=[r['id'] for r in req if (r['area']=='FUNC')==(fig=='functional-requirements')]
 assert all(x in svg for x in expected)
D=Decimal
def money(x):return D(x).quantize(D('.01'),rounding=ROUND_HALF_UP)
def end(b,i,p,e,l,days):return money(D(b)+D(i)-D(p)-D(e)/D(l)*D(days))
def next_month(e,l,n,p):return money(D(e)/D(l)*D(n)+D(p))
def full_months(start,now):
 year,month=now.year,now.month;result=[]
 for _ in range(3):
  month-=1
  if month==0:year-=1;month=12
  first=date(year,month,1)
  if first>=start:result.append((year,month))
 return list(reversed(result))
checks={
 'F01':end(1000,0,100,920,92,5)==D('850.00'),
 'F02':next_month(920,92,30,100)==D('400.00'),
 'F03':full_months(date(2026,10,15),date(2026,10,26))==[],
 'F04':end(1000,0,100,0,30,5)==D('900.00'),
 'F05':end(1000,0,100,920,92,0)==D('900.00'),
 'F06':end(100,0,100,920,92,5)==D('-50.00'),
 'F07':end(1000,0,100,920,92,5)==end(900,0,0,920,92,5)==D('850.00'),
 'F09':end(1000,0,0,300,61,5)==D('975.41'),
 'F10':min(31,calendar.monthrange(2028,2)[1])==29,
 'F11':min(31,calendar.monthrange(2027,2)[1])==28,
 'F12':end(1000,0,100,920,92,5)==D('850.00') and end(1000,0,0,920,92,5)==D('950.00'),
 'F13':full_months(date(2026,7,15),date(2026,10,26))==[(2026,8),(2026,9)],
 'F14':full_months(date(2026,6,1),date(2026,10,26))==[(2026,7),(2026,8),(2026,9)],
 'F15':end(900,0,0,920,92,5)==end(1000,0,100,920,92,5),
 'F16':end(990,0,100,930,92,5)==D('839.46') and next_month(930,92,30,100)==D('403.26'),
 'F17':end(1000,0,0,920,92,5)==D('950.00') and next_month(920,92,30,200)==D('500.00'),
}
assert all(checks.values()),checks
print(f'PASS: {len(ids)} requirements match SRS, SysML, acceptance scenarios and matrix; links exist.')
print(f'PASS: {len(checks)} arithmetic/calendar checks. Duplicate-request behavior F08 remains planned for application tests.')
