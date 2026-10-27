#!/usr/bin/env python3
"""Editable SVG diagrams and redesigned interface mockups for the lab report.

Run with Python 3, then build_report_png.cjs (Node.js and sharp). These are report illustrations, not
screenshots of an implemented application or direct SysML Pilot exports.
"""
import html
import json
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
OUT = ROOT / 'reports/lab3/figures'
REQ = json.loads((ROOT / 'srs/data/requirements.json').read_text())
INK, MUTED, LINE, BLUE = '#1f2937', '#6b7280', '#d4dce6', '#0284c7'


class SVG:
    def __init__(self, w, h):
        self.w, self.h = w, h
        self.s = [f'<svg xmlns="http://www.w3.org/2000/svg" width="{w*2}" height="{h*2}" viewBox="0 0 {w} {h}">',
                  '<defs><marker id="arrow" markerWidth="12" markerHeight="10" refX="10" refY="5" orient="auto"><path d="M1,1 L10,5 L1,9" fill="none" stroke="#64748b" stroke-width="1.4"/></marker></defs>']
        self.rect(0, 0, w, h, '#ffffff', 'none', 0)

    def rect(self, x, y, w, h, fill='#ffffff', stroke=LINE, radius=7):
        self.s.append(f'<rect x="{x}" y="{y}" width="{w}" height="{h}" rx="{radius}" fill="{fill}" stroke="{stroke}"/>')

    def text(self, x, y, s, size=16, color=INK, bold=False, anchor='start'):
        self.s.append(f'<text x="{x}" y="{y}" font-family="Arial, DejaVu Sans, sans-serif" font-size="{size}" fill="{color}" font-weight="{700 if bold else 400}" text-anchor="{anchor}">{html.escape(str(s))}</text>')

    def lines(self, x, y, items, size=16, color=INK, bold=False, step=None, anchor='start'):
        for i, s in enumerate(items):
            self.text(x, y+i*(step or size*1.35), s, size, color, bold, anchor)

    def line(self, x1, y1, x2, y2, color=LINE, width=1.3, dash='', arrow=False):
        self.s.append(f'<line x1="{x1}" y1="{y1}" x2="{x2}" y2="{y2}" stroke="{color}" stroke-width="{width}" stroke-dasharray="{dash}"'+(' marker-end="url(#arrow)"' if arrow else '')+'/>')

    def path(self, d, color=LINE, width=1.3):
        self.s.append(f'<path d="{d}" fill="none" stroke="{color}" stroke-width="{width}"/>')

    def button(self, x, y, w, label, primary=False, color=None):
        self.rect(x, y, w, 34, color or (BLUE if primary else '#ffffff'), color or (BLUE if primary else '#cbd5e1'), 6)
        self.text(x+w/2, y+23, label, 14, '#ffffff' if primary or color else INK, True, 'middle')

    def field(self, x, y, w, label, value, error=False):
        self.text(x, y, label, 14, MUTED)
        self.rect(x, y+9, w, 36, '#ffffff', '#ef4444' if error else '#cbd5e1', 5)
        self.text(x+12, y+33, value, 15)

    def chip(self, x, y, label, kind='blue', w=None):
        fg, bg = {'blue':('#1d4ed8','#dbeafe'),'green':('#047857','#d1fae5'),'red':('#b91c1c','#fee2e2'),'gray':('#4b5563','#e5e7eb'),'amber':('#92400e','#fef3c7')}[kind]
        w = w or len(label)*8+22
        self.rect(x,y,w,25,bg,'none',13)
        self.text(x+11,y+18,label,13,fg,True)

    def save(self, name):
        OUT.mkdir(parents=True, exist_ok=True)
        data = ('\n'.join(self.s)+'</svg>').encode()
        (OUT / f'{name}.svg').write_bytes(data)


# Section 3.5 uses a distinct warm, green design. UML assets retain their own style.
UI_INK, UI_MUTED, UI_LINE = '#203d32', '#63736a', '#dce3d9'
GREEN, SAGE, ORANGE = '#245c48', '#9bb693', '#d59565'


class UI(SVG):
    def rect(self, x, y, w, h, fill='#ffffff', stroke=UI_LINE, radius=14):
        super().rect(x,y,w,h,fill,stroke,radius)

    def text(self, x, y, value, size=16, color=UI_INK, bold=False, anchor='start'):
        super().text(x,y,value,size,color,bold,anchor)

    def lines(self, x, y, items, size=16, color=UI_MUTED, bold=False, step=None, anchor='start'):
        super().lines(x,y,items,size,color,bold,step,anchor)

    def button(self, x, y, w, label, primary=False, color=None):
        fill = color or (GREEN if primary else '#ffffff')
        self.rect(x,y,w,38,fill,fill if primary or color else UI_LINE,19)
        self.text(x+w/2,y+25,label,15,'#ffffff' if primary or color else GREEN,True,'middle')

    def field(self, x, y, w, label, value, error=False):
        self.text(x,y,label,14,UI_MUTED)
        self.rect(x,y+10,w,38,'#ffffff','#b4513b' if error else '#c9d4c5',9)
        self.text(x+12,y+35,value,16)

    def chip(self, x, y, label, kind='green', w=None):
        fg,bg={'green':(GREEN,'#e8f0e3'),'amber':('#99592f','#f8ecde'),
               'red':('#9c442e','#fae8e1'),'gray':('#63736a','#edf0eb')}[kind]
        w=w or len(label)*8+22
        self.rect(x,y,w,26,bg,'none',13)
        self.text(x+w/2,y+18,label,13,fg,True,'middle')

    def ring(self, cx, cy, r, items, width=23):
        import math
        circumference=2*math.pi*r
        offset=0
        for share,color in items:
            length=circumference*share
            self.s.append(f'<circle cx="{cx}" cy="{cy}" r="{r}" fill="none" stroke="{color}" stroke-width="{width}" stroke-dasharray="{length-3} {circumference-length+3}" stroke-dashoffset="{-offset}" transform="rotate(-90 {cx} {cy})"/>')
            offset+=length


def ui(title, active, h=720, subtitle=''):
    s=UI(1100,h)
    s.rect(0,0,1100,h,'#f5f5ef','none',0)
    s.rect(0,0,1100,96,'#ffffff','none',0)
    s.rect(28,14,31,31,GREEN,'none',10)
    for x,y,hh in [(36,29,9),(43,24,14),(50,20,18)]:s.rect(x,y,3,hh,'#f5f5ef','none',1)
    s.text(70,38,'баланс',26,GREEN,True)
    s.text(191,36,'Личный бюджет',15,UI_MUTED)
    s.chip(813,19,'BYN',w=58)
    s.text(889,36,'Мой профиль',15)
    s.text(1068,36,'Выйти',14,UI_MUTED,False,'end')
    tabs=[('Обзор',28,100),('Операции',144,134),('Категории и лимиты',294,218),('Плановые платежи',528,214),('Аналитика и прогнозы',758,271)]
    for i,(label,x,w) in enumerate(tabs):
        if i==active:s.rect(x,58,w,30,'#e8f0e3','none',15)
        s.text(x+w/2,79,label,16,GREEN if i==active else UI_MUTED,i==active,'middle')
    s.text(28,132,title,27,UI_INK,True)
    s.text(28,157,subtitle,15,UI_MUTED)
    return s


def overview():
    s=ui('Мой бюджет',0,640,'Октябрь 2026 · на конец дня 26 октября')
    s.button(878,114,194,'+ Добавить операцию',True)
    s.rect(28,177,366,229,GREEN,'none',20)
    s.text(52,211,'ТЕКУЩИЙ ОСТАТОК',14,'#cfe0c9',True)
    s.text(52,262,'1 000,00',44,'#ffffff',True)
    s.text(52,291,'BYN',19,'#dfeadb')
    s.line(52,313,370,313,'#527663')
    s.text(52,345,'Начальный остаток',16,'#dfeadb')
    s.text(370,345,'500,00 BYN',18,'#ffffff',True,'end')
    s.text(52,382,'Один бюджет · ручной учёт',14,'#cfe0c9')
    for x,label,value,small in [(414,'Доходы за месяц','1 500,00','Поступления в октябре'),(753,'Расходы за месяц','1 000,00','Фактические траты')]:
        s.rect(x,177,319,109)
        s.text(x+20,207,label,16,UI_MUTED)
        s.text(x+20,244,value,29,GREEN,True)
        s.text(x+186,244,'BYN',16,UI_MUTED)
        s.text(x+20,269,small,13,UI_MUTED)
    s.rect(414,302,658,104,'#e8efdf','none',16)
    s.text(436,331,'Остаток к 31 октября',17,GREEN,True)
    s.text(436,375,'850,00 BYN',31,GREEN,True)
    s.lines(760,340,['С учётом планов','и повседневных трат'],16)
    s.text(760,387,'Открыть прогноз  →',15,GREEN,True)
    s.rect(28,424,366,186)
    s.text(48,455,'Ближайший платёж',19,UI_INK,True)
    s.rect(48,475,66,62,'#f0f2e9','none',12)
    s.text(81,501,'30',24,GREEN,True,'middle');s.text(81,523,'ОКТ',12,UI_MUTED,True,'middle')
    s.text(132,496,'Интернет',19,UI_INK,True);s.text(132,524,'100,00 BYN',20,GREEN,True)
    s.text(48,568,'Ожидается подтверждение',15,UI_MUTED)
    s.text(48,591,'План ещё не изменяет остаток.',14,UI_MUTED)
    s.rect(414,424,365,186)
    s.text(434,455,'Структура расходов',19,UI_INK,True)
    s.ring(478,531,40,[(.45,GREEN),(.15,ORANGE),(.4,SAGE)],17)
    for i,(name,amount,c) in enumerate([('Продукты','450',GREEN),('Транспорт','150',ORANGE),('Жильё','400',SAGE)]):
        y=494+i*34
        s.rect(538,y-10,9,9,c,'none',3);s.text(557,y,name,15);s.text(758,y,amount,16,UI_INK,True,'end')
    s.text(758,593,'BYN',13,UI_MUTED,False,'end')
    s.rect(799,424,273,186,'#fbebdf','#e9cbb5',16)
    s.text(819,456,'Превышен лимит',19,'#95522c',True)
    s.text(819,488,'Транспорт',16,'#95522c')
    s.text(819,527,'+30,00 BYN',28,'#95522c',True)
    s.lines(819,562,['Потрачено 150 из 120 BYN.','Операции доступны.'],14,'#95522c')
    s.save('04-overview')


def operations():
    s=ui('Операции',1,780,'Все доходы и расходы в одном списке')
    s.button(878,114,194,'+ Добавить операцию',True)
    s.rect(28,177,638,62)
    s.button(42,189,224,'01.10–26.10.2026 ▾')
    s.button(278,189,164,'Все типы ▾')
    s.button(454,189,198,'Все категории ▾')
    s.rect(28,255,638,338)
    s.text(48,286,'Последние операции',20,UI_INK,True)
    s.text(646,286,'Сумма, BYN',14,UI_MUTED,False,'end')
    rows=[('26 окт','Продукты','Покупка продуктов','−45,50'),('24 окт','Транспорт','Поездка','−15,00'),('20 окт','Зарплата','Поступление зарплаты','+1 500,00'),('15 окт','Жильё','Оплата жилья','−400,00')]
    for i,(date,cat,desc,value) in enumerate(rows):
        y=306+i*61
        if i==0:s.rect(38,y,618,58,'#eef3e8','none',10)
        s.text(49,y+28,date,15,UI_MUTED)
        s.text(127,y+23,cat,17,UI_INK,True);s.text(127,y+45,desc,14,UI_MUTED)
        s.text(644,y+29,value,19,GREEN if value[0]=='+' else UI_INK,True,'end')
        if i<3:s.line(49,y+60,645,y+60,UI_LINE)
    s.text(48,576,'4 последние записи',14,UI_MUTED);s.text(645,576,'‹  1  2  3  ›',16,GREEN,True,'end')
    s.rect(686,177,386,416)
    s.chip(707,193,'Расход',w=93);s.text(1049,212,'26 октября',14,UI_MUTED,False,'end')
    s.text(707,249,'Редактирование операции',22,UI_INK,True)
    s.field(707,279,160,'Дата','26.10.2026');s.field(884,279,167,'Сумма, BYN','0,00',True)
    s.text(884,349,'Введите сумму > 0',13,'#9c442e')
    s.field(707,376,344,'Категория','Продукты ▾')
    s.field(707,451,344,'Комментарий','Покупка продуктов')
    s.button(707,518,159,'Сохранить',True);s.button(882,518,169,'Отмена')
    s.text(707,580,'Удалить операцию',14,'#a24b35',True)
    s.text(28,619,'ДРУГИЕ СОСТОЯНИЯ ЭКРАНА',12,UI_MUTED,True)
    s.rect(28,633,474,122,'#eaf0e3','none',14)
    s.text(48,663,'Нет операций за выбранный период',19,UI_INK,True)
    s.text(48,690,'Добавьте первую запись или измените фильтры.',15,UI_MUTED)
    s.text(48,731,'Сбросить фильтры  →',15,GREEN,True)
    s.rect(522,633,550,122,'#fffaf6','#ead6c6',14)
    s.text(542,663,'Удалить расход 45,50 BYN?',19,UI_INK,True)
    s.text(542,690,'Остаток увеличится на 45,50 BYN.',15,UI_MUTED)
    s.button(542,705,131,'Удалить',color='#a24b35');s.button(686,705,122,'Отмена')
    s.text(1049,732,'Подтверждение',13,UI_MUTED,False,'end')
    s.save('05-operations')


def categories():
    s=ui('Категории и лимиты',2,720,'План расходов на октябрь 2026')
    s.button(878,114,194,'+ Новая категория',True)
    s.rect(28,177,657,515)
    s.text(48,207,'Расходы по лимитам',20,UI_INK,True)
    s.text(665,207,'BYN',14,UI_MUTED,False,'end')
    s.chip(48,223,'В пределах лимита',w=174)
    for i,(name,amount,limit,pct) in enumerate([('Продукты',450,500,.9),('Жильё',400,450,400/450)]):
        y=279+i*99
        s.text(48,y,name,19,UI_INK,True);s.text(665,y,f'{amount},00 / {limit},00',17,UI_INK,True,'end')
        s.rect(48,y+17,617,10,'#e9ede5','none',5);s.rect(48,y+17,617*pct,10,GREEN,'none',5)
        s.text(48,y+50,'Осталось 50,00 BYN',15,UI_MUTED)
    s.chip(48,447,'Лимит превышен',kind='amber',w=159)
    s.text(48,505,'Транспорт',19,UI_INK,True);s.text(665,505,'150,00 / 120,00',17,'#99592f',True,'end')
    s.rect(48,522,617,10,ORANGE,'none',5)
    s.text(48,555,'Перерасход 30,00 BYN',15,'#99592f',True)
    s.line(48,580,665,580,UI_LINE)
    s.chip(48,598,'Без лимита',kind='gray',w=120)
    s.text(48,656,'Зарплата',19,UI_INK,True);s.text(665,656,'Доход · лимит не применяется',15,UI_MUTED,False,'end')
    s.rect(705,177,367,422)
    s.text(725,210,'Настройка категории',22,UI_INK,True)
    s.field(725,247,327,'Название','Транспорт')
    s.field(725,327,327,'Тип','Расход ▾')
    s.field(725,407,327,'Лимит на месяц, BYN','120,00')
    s.button(725,480,327,'Сохранить изменения',True)
    s.text(725,553,'Архивировать категорию',15,'#a24b35',True)
    s.text(725,580,'История операций сохранится.',14,UI_MUTED)
    s.rect(705,617,367,75,'#eaf0e3','none',14)
    s.lines(725,647,['Архивные категории недоступны','для новых операций.'],16)
    s.save('06-categories')


def payments():
    s=ui('Плановые платежи',3,740,'Планируйте регулярные поступления и расходы')
    s.button(878,114,194,'+ Новое правило',True)
    s.rect(28,177,1044,79,'#eaf0e3','none',14)
    s.text(48,207,'ИНТЕРНЕТ',13,GREEN,True);s.text(48,237,'Ежемесячный расход',20,UI_INK,True)
    s.text(449,207,'ПОВТОРЕНИЕ',13,UI_MUTED,True);s.text(449,237,'Каждого 30 числа',20,UI_INK,True)
    s.text(1052,207,'СУММА',13,UI_MUTED,True,'end');s.text(1052,237,'100,00 BYN',24,GREEN,True,'end')
    s.rect(28,274,614,438)
    s.text(48,307,'Расписание',21,UI_INK,True);s.text(621,307,'Октябрь — ноябрь',15,UI_MUTED,False,'end')
    rows=[('30','НОЯ','Интернет','100,00 BYN','Ожидается','green'),('30','ОКТ','Интернет','100,00 BYN','Ожидается','green'),('20','ОКТ','Зарплата','1 500,00 BYN','Исполнено','green'),('05','ОКТ','Подписка','20,00 BYN','Пропущено','gray')]
    for i,(day,mon,name,amount,status,color) in enumerate(rows):
        y=325+i*76
        if i==1:s.rect(39,y,592,71,'#eef3e8','none',12)
        s.rect(49,y+10,53,51,'#f0f2e9','none',10)
        s.text(75,y+33,day,24,GREEN,True,'middle');s.text(75,y+50,mon,11,UI_MUTED,True,'middle')
        s.text(119,y+30,name,18,UI_INK,True);s.text(119,y+55,amount,16,UI_MUTED)
        s.chip(490,y+23,status,color,126)
    s.lines(48,663,['В коротком месяце — последний день.','Перенос действует только на выбранную дату.'],15)
    s.rect(662,274,410,438)
    s.chip(683,292,'Ожидается',w=121)
    s.text(683,348,'Интернет · 30 октября',23,UI_INK,True)
    s.field(683,382,176,'Фактическая дата','26.10.2026');s.field(876,382,175,'Сумма, BYN','100,00')
    s.button(683,453,368,'Создать операцию',True)
    s.button(683,504,368,'Связать существующую')
    s.button(683,555,176,'Пропустить');s.button(876,555,175,'Перенести')
    s.lines(683,632,['После исполнения платёж исключается','из будущих сумм. Повторное подтверждение','не создаёт вторую операцию.'],15)
    s.save('07-payments')


def analytics():
    s=ui('Аналитика и прогнозы',4,720,'Расчёт на конец дня 26.10.2026 · все суммы в BYN')
    s.button(878,114,194,'Октябрь 2026 ▾')
    s.rect(28,177,512,152,GREEN,'none',17)
    s.text(49,208,'ОСТАТОК НА 31 ОКТЯБРЯ',14,'#d0dfc9',True)
    s.text(49,254,'850,00 BYN',37,'#ffffff',True)
    s.text(49,287,'1 000 + 0 − 100 − 10 × 5 = 850',17,'#e2ebdc')
    s.text(49,313,'5 дней · будущий платёж 100 BYN',14,'#d0dfc9')
    s.rect(560,177,512,152,'#e6eddb','none',17)
    s.text(581,208,'РАСХОДЫ В НОЯБРЕ',14,GREEN,True)
    s.text(581,254,'400,00 BYN',37,GREEN,True)
    s.text(581,287,'10 × 30 + 100 = 400',17,GREEN)
    s.text(581,313,'30 дней · плановые расходы 100 BYN',14,UI_MUTED)
    s.rect(28,347,512,214)
    s.text(49,379,'Структура расходов',21,UI_INK,True)
    s.ring(123,464,57,[(.45,GREEN),(.15,ORANGE),(.4,SAGE)],23)
    s.text(123,460,'1 000',22,GREEN,True,'middle');s.text(123,481,'BYN',13,UI_MUTED,False,'middle')
    for i,(name,value,c) in enumerate([('Продукты','450 · 45 %',GREEN),('Транспорт','150 · 15 %',ORANGE),('Жильё','400 · 40 %',SAGE)]):
        y=423+i*43
        s.rect(224,y-11,9,9,c,'none',3);s.text(243,y,name,16);s.text(520,y,value,16,UI_INK,True,'end')
    s.rect(560,347,512,214)
    s.text(581,379,'Сравнение расходов',21,UI_INK,True)
    s.text(581,406,'01–26 октября и 01–26 сентября 2026',15,UI_MUTED)
    s.text(581,443,'Октябрь',15);s.rect(684,426,268,23,GREEN,'none',5);s.text(1051,444,'1 000',16,UI_INK,True,'end')
    s.text(581,482,'Сентябрь',15);s.rect(684,465,241,23,SAGE,'none',5);s.text(1051,483,'900',16,UI_INK,True,'end')
    s.text(581,534,'+100,00 BYN (+11,11 %) · неполные месяцы',15,UI_MUTED)
    s.text(28,591,'Как получен прогноз',18,UI_INK,True)
    s.text(28,618,'Июль–сентябрь: 920 / 92 = 10 BYN/день (нерегулярные траты). Проведённые платежи не дублируются.',16,UI_MUTED)
    s.rect(28,637,1044,63,'#eaf0e3','none',14)
    s.text(48,662,'Нет полного месяца истории → «Недостаточно данных для прогноза».',16,GREEN,True)
    s.text(48,686,'Прогноз оценивает будущие суммы при сохранении уровня трат и не является гарантией.',15,UI_MUTED)
    s.save('08-analytics')


def actor(s,x,y,label):
    s.s.append(f'<circle cx="{x}" cy="{y}" r="16" fill="white" stroke="#334155" stroke-width="2"/>')
    for a,b,c,d in [(x,y+16,x,y+65),(x-28,y+36,x+28,y+36),(x,y+65,x-24,y+100),(x,y+65,x+24,y+100)]:s.line(a,b,c,d,'#334155',2)
    s.text(x,y+128,label,19,INK,False,'middle')


def usecases():
    s=SVG(1100,1000)
    s.rect(300,25,770,950,'#f8fafc','#334155',22)
    s.text(685,69,'Баланс — личный бюджет',26,INK,True,'middle')
    actor(s,125,120,'Гость');actor(s,125,525,'Владелец бюджета')
    cases=[(710,140+i*100,[name]) for i,name in enumerate(['Регистрация','Вход и выход','Начало учёта','Обзор бюджета','Операции и фильтры','Категории и лимиты','Плановые платежи','Аналитика и прогнозы'])]
    for i,(x,y,labels) in enumerate(cases):
        if i<2:s.line(153,156,x-270,y,'#64748b',1.7)
        else:s.line(153,561,x-270,y,'#64748b',1.7)
    for x,y,labels in cases:
        s.s.append(f'<ellipse cx="{x}" cy="{y}" rx="270" ry="38" fill="white" stroke="#334155" stroke-width="2"/>')
        s.lines(x,y+6,labels,19,anchor='middle')
    s.lines(685,922,['Владелец работает только с собственным бюджетом.',
                     'Фактические операции и планы учитываются раздельно.'],18,MUTED,anchor='middle')
    s.save('01-use-cases')


def tree(name, groups, title):
    # Requirement labels come from the same catalogue as the SRS and SysML.
    total=sum(len(items) for _,items in groups)
    h=total*72+70
    s=SVG(1100,h)
    center=h/2
    s.rect(20,center-68,228,136,'#e2e8f0','#64748b',0)
    s.lines(134,center-38,['«requirement»','BalanceSpec',title[0],title[1]],19,INK,False,29,'middle')
    y=28
    for label,items in groups:
        cy=y+(len(items)*72-12)/2
        s.line(248,center,302,cy,'#64748b',1.5)
        s.rect(302,cy-48,240,96,'#edf2f7','#64748b',0)
        s.lines(422,cy-12,['«requirement»',label],18,INK,False,29,'middle')
        for r in items:
            s.line(542,cy,600,y+29,'#64748b',1.4)
            s.rect(600,y,480,58,'#ffffff','#64748b',0)
            s.text(615,y+23,r['id'],18,INK,True)
            s.text(615,y+47,r['title'],18)
            y+=72
    s.save(name)


def requirements():
    functional=[r for r in REQ if r['area']=='FUNC']
    groups=[('Доступ',functional[:2]),('Бюджет и операции',functional[2:8]),('Аналитика и лимиты',functional[8:11]),('Плановые платежи',functional[11:13]+functional[17:]),('Прогнозирование',functional[13:17])]
    tree('02-functional',groups,['Функциональные','требования'])
    n=[r for r in REQ if r['area']!='FUNC']
    tree('03-quality',[(g,[r for r in n if r['group']==g]) for g in ['Интерфейсы','Качество','Безопасность','Ограничения']],['Нефункциональные','требования'])


def domain():
    s=SVG(1200,1020)
    def box(x,y,w,title,fields):
        h=45+len(fields)*29
        s.rect(x,y,w,h,'white','#334155',0);s.rect(x,y,w,39,'#e2e8f0','#334155',0)
        s.text(x+w/2,y+27,title,22,INK,True,'middle')
        s.lines(x+14,y+66,fields,18,step=29)
    # Routing avoids class interiors; a bridge denotes a crossing without a junction.
    for coords in [(310,100,460,100),(740,100,910,100),(600,186,600,650),(1080,215,1080,360),(180,608,180,680),(330,790,460,790)]:s.line(*coords,'#64748b',1.8)
    s.path('M 500,186 V 260 H 160 V 360','#64748b',1.8)
    s.path('M 910,165 H 820 V 570 H 780 V 720 H 760','#64748b',1.8)
    s.path('M 910,195 H 850 V 300 H 610 Q 600,285 590,300 H 290 V 360','#64748b',1.8)
    for x,y,t in [(326,90,'1'),(430,90,'1'),(758,90,'1'),(863,90,'0..*'),(514,213,'1'),(173,347,'0..*'),(614,212,'1'),(614,636,'0..*'),(1092,244,'1'),(1092,347,'0..*'),(838,157,'1'),(786,706,'0..*'),(857,215,'1'),(300,347,'0..*'),(191,633,'0..1'),(191,667,'0..1'),(342,780,'0..*'),(435,780,'1')]:s.text(x,y,t,17,MUTED)
    box(30,25,280,'User',['id: UUID','email: String','password_hash: String'])
    box(460,25,280,'Budget',['id: UUID','user_id: UUID','start_date: Date','initial_balance: Money'])
    box(910,25,260,'Category',['id: UUID','budget_id: UUID','type: Income | Expense','name: String','archived: Boolean'])
    box(30,360,300,'Operation',['id: UUID','budget_id: UUID','type: Income | Expense','category_id: UUID','date: Date','amount: Money','comment: String?'])
    box(30,680,300,'Occurrence',['id: UUID','rule_id: UUID','due_date: Date','amount: Money','status: Status','operation_id: UUID?'])
    box(460,650,300,'RecurringRule',['id: UUID','budget_id: UUID','type: Income | Expense','category_id: UUID','amount: Money','day: Integer (1–31)','start: Date / end: Date?'])
    box(870,360,300,'CategoryLimit',['budget_id: UUID','category_id: UUID','month: YearMonth','amount: Money'])
    s.lines(600,968,['Money — сумма в BYN с точностью до копейки. Прогноз — производные данные.',
                      'Status: ожидается, исполнено, пропущено. Просрочка определяется по дате.'],17,MUTED,anchor='middle')
    s.save('09-domain')


if __name__=='__main__':
    usecases();requirements();overview();operations();categories();payments();analytics();domain()
    print(f'Created 9 SVG sources in {OUT}')
