"""Validate documentation consistency and specification arithmetic, not application code."""
from pathlib import Path
from decimal import Decimal, ROUND_HALF_UP
from datetime import date
import json,re,calendar,hashlib
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
def heading_anchors(text):
 return {re.sub(r'[^\w\- ]','',h.lower()).replace(' ','-') for h in re.findall(r'^#{1,6} (.+)$',text,re.M)}
for file in [S/'srs.md',S/'sysml/README.md',R/'README.md',S/'verification/document-checks.md',R/'reports/lab3/figures/README.md']:
 for link in re.findall(r'\]\(([^)]+)\)',file.read_text()):
  if link.startswith(('https://','http://')):continue
  path,separator,anchor=link.partition('#')
  target=file.parent/path if path else file
  assert target.exists(),(file,link)
  if separator and target.suffix=='.md':assert anchor in heading_anchors(target.read_text()),(file,link)
views=json.loads((S/'sysml/tools/views.json').read_text())
assert len(views)==8
manifest=json.loads((S/'sysml/tools/render-manifest.json').read_text())
for filename,digest in manifest['sources'].items():
 assert hashlib.sha256((S/'sysml'/filename).read_bytes()).hexdigest()==digest,filename
for filename,digest in manifest['tooling'].items():
 assert hashlib.sha256((S/'sysml/tools'/filename).read_bytes()).hexdigest()==digest,filename
for v in views:
 filename=v['name']+'.png'
 assert (S/'images'/filename).read_bytes().startswith(b'\x89PNG\r\n\x1a\n')
 assert hashlib.sha256((S/'images'/filename).read_bytes()).hexdigest()==manifest['images'][filename],filename
 assert v['file'] in manifest['sources']
trace=(S/'sysml/use-cases.sysml').read_text()
for r in req:
 if r['area']=='FUNC':assert '::'+r['id'].replace('-','_')+';' in trace
assert not (S/'mockups').exists() and not (S/'img').exists()
# The report uses a separate illustrated set in the supplied DOCX style.
from zipfile import ZipFile
from xml.etree import ElementTree as ET
figures=R/'reports/lab3/figures'
names=['01-use-cases','02-functional','03-quality','04-overview','05-operations','06-categories','07-payments','08-analytics','09-domain']
assert {p.stem for p in figures.glob('*.svg')}==set(names)
assert {p.stem for p in figures.glob('*.png')}==set(names)
for r in req:
 tree=figures/('02-functional.svg' if r['area']=='FUNC' else '03-quality.svg')
 assert tree.read_text().count(r['id'])==1,r['id']
with ZipFile(R/'reports/lab3/Отчёт_ЛР3_Баланс.docx') as doc:
 xml=ET.fromstring(doc.read('word/document.xml'))
 ns={'w':'http://schemas.openxmlformats.org/wordprocessingml/2006/main'}
 text=' '.join(xml.itertext())
 assert len(xml.findall('.//w:tbl',ns))==6
 assert all(f'Вопрос {i}.' in text for i in range(1,7))
 assert not any(t in text for t in ['FreelanceDesk','Могильный','Апанас','Новицкий','srs/mockups'])
 for i,name in enumerate(names,1):
  assert doc.read(f'word/media/image{i}.png')==(figures/f'{name}.png').read_bytes(),name
print('PASS: 9 report figures match DOCX; 37 tree IDs, 6 tables and 6 questions checked.')
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
print(f'PASS: {len(checks)} arithmetic/calendar checks. Application behavior F08 is tested separately by go test ./....')

print('PASS: 8 Pilot exports match source hashes; navigation links and functional traceability checked.')
