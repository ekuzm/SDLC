// Optional browser smoke test. Requires playwright and a running EMPTY demo instance.
// Example: BALANCE_URL=http://127.0.0.1:18080 node tools/smoke.cjs
const {chromium} = require('playwright');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
(async()=>{
 const engine=process.env.BALANCE_BROWSER==='firefox'?require('playwright').firefox:chromium;
 const browser=await engine.launch({headless:true});
 try{
 const context=await browser.newContext({viewport:{width:1440,height:1000},locale:'ru-RU'});
 const page=await context.newPage();const base=process.env.BALANCE_URL||'http://127.0.0.1:8080';
 const out=process.env.BALANCE_SCREENSHOTS||'/tmp/balance-smoke';fs.mkdirSync(out,{recursive:true});
 const email=`demo-${Date.now()}@example.by`;const password=`Test-only-${Date.now()}-budget`;
 const errors=[];page.on('pageerror',error=>errors.push(error.message));
 const submit=async(form,button)=>{await Promise.all([page.waitForURL(url=>url.searchParams.get('saved')==='1'),form.getByRole('button',{name:button,exact:true}).click()]);};
 await page.goto(base+'/register');await page.getByLabel('Email',{exact:true}).fill(email);await page.getByLabel('Пароль',{exact:true}).fill(password);
 await Promise.all([page.waitForURL('**/overview'),page.getByRole('button',{name:'Зарегистрироваться',exact:true}).click()]);
 await page.goto(base+'/settings');let form=page.locator('form[action="/budget"]');await form.locator('[name=start]').fill('2026-07-01');await form.locator('[name=initial]').fill('1920');await submit(form,'Сохранить бюджет');
 const today=new Date().toLocaleDateString('en-CA',{timeZone:'Europe/Minsk'});
 const add=async(amount,date,comment,kind='expense',category='Продукты')=>{await page.goto(base+'/operations');const f=page.locator('form[action="/operation"]');await f.locator('[name=kind]').selectOption(kind);await f.locator('[name=category]').selectOption({label:category+' · '+(kind==='income'?'Доход':'Расход')});await f.locator('[name=date]').fill(date);await f.locator('[name=amount]').fill(amount);await f.locator('[name=comment]').fill(comment);await submit(f,'Сохранить');};
 await add('920','2026-07-15','История: нерегулярные расходы');
 // A rejected operation must preserve the comment and leave the balance unchanged.
 await page.goto(base+'/operations');form=page.locator('form[action="/operation"]');await form.locator('[name=category]').selectOption({label:'Продукты · Расход'});await form.locator('[name=amount]').fill('0');await form.locator('[name=comment]').fill('Проверка ошибки');
 await Promise.all([page.waitForResponse(r=>r.url().endsWith('/operation')&&r.request().method()==='POST'),form.getByRole('button',{name:'Сохранить',exact:true}).click()]);
 assert.equal(await page.locator('[name=comment]').inputValue(),'Проверка ошибки');assert.match(await page.getByRole('alert').innerText(),/0,01/);
 await add('20',today,'Проверка редактирования');await page.getByRole('link',{name:'Изменить',exact:true}).first().click();form=page.locator('form[action="/operation"]');await form.locator('[name=amount]').fill('25');await submit(form,'Сохранить');assert.match(await page.locator('table').innerText(),/25,00/);
 await page.getByRole('link',{name:'Удалить',exact:true}).first().click();assert.equal(await page.locator('form[action="/operation/delete"]').count(),1);await page.getByRole('link',{name:'Отмена',exact:true}).click();assert.match(await page.locator('table').innerText(),/25,00/);
 await page.getByRole('link',{name:'Удалить',exact:true}).first().click();form=page.locator('form[action="/operation/delete"]');await form.getByRole('checkbox').check();await submit(form,'Удалить операцию');
 await page.goto(base+'/categories');form=page.locator('article.category-row').filter({has:page.getByRole('heading',{name:'Продукты',exact:true})}).locator('form[action="/limit"]');await form.locator('[name=amount]').fill('500');await submit(form,'Сохранить');
 await page.goto(base+'/payments');form=page.locator('form[action="/rule"]');await form.locator('[name=name]').fill('Интернет');await form.locator('[name=category]').selectOption({label:'Другое · Расход'});await form.locator('[name=amount]').fill('100');await form.locator('[name=day]').fill('30');await form.locator('[name=start]').fill('2026-10-01');await submit(form,'Сохранить правило');
 await page.goto(base+'/analytics');assert.match(await page.locator('main').innerText(),/400,00/);assert.match(await page.locator('main').innerText(),/920,00/);
 // Confirm a regular payment and replay it via HTTP; only one actual operation may appear.
 await page.goto(base+'/payments');let details=page.locator('details.payment').first();await details.locator('summary').click();let payment=details.locator('form');const csrf=await payment.locator('[name=csrf]').inputValue();const id=await payment.locator('[name=id]').inputValue();await payment.locator('[name=date]').fill(today);await submit(payment,'Создать операцию');
 const replay=await context.request.post(base+'/payment',{form:{csrf,id,action:'complete',date:today,amount:'100'},maxRedirects:0});assert.equal(replay.status(),303);
 await page.goto(base+'/operations');assert.equal(await page.locator('tbody tr').count(),2);
 // Deleting the linked expense restores pending status without inflating the forecast.
 await page.getByRole('link',{name:'Удалить',exact:true}).first().click();form=page.locator('form[action="/operation/delete"]');await form.getByRole('checkbox').check();await submit(form,'Удалить операцию');
 await page.goto(base+'/payments');assert.match(await page.locator('details.payment').first().innerText(),/Ожидается/);
 // Link an operation entered manually, then delete it and reopen the plan.
 await add('100',today,'Интернет (введён вручную)','expense','Другое');
 await page.goto(base+'/payments');details=page.locator('details.payment').first();await details.locator('summary').click();payment=details.locator('form');
 const options=await payment.locator('select[name=operation] option').allTextContents();const label=options.find(s=>s.includes('Другое')&&s.includes('100,00 BYN'));assert.ok(label,'manual operation offered for linking');
 await payment.locator('select[name=operation]').selectOption({label});await submit(payment,'Связать существующую операцию');
 await page.goto(base+'/operations');assert.equal(await page.locator('tbody tr').count(),2);await page.getByRole('link',{name:'Удалить',exact:true}).first().click();form=page.locator('form[action="/operation/delete"]');await form.getByRole('checkbox').check();await submit(form,'Удалить операцию');
 // Add an illustrative current-month budget for the screenshots after assertions.
 await page.goto(base+'/settings');form=page.locator('form[action="/budget"]');await form.locator('[name=initial]').fill('1420');await form.getByRole('checkbox').check();await submit(form,'Сохранить бюджет');
 await add('1500','2026-10-20','Зарплата за октябрь','income','Зарплата');
 await add('450','2026-10-26','Продукты за месяц');
 await add('150','2026-10-24','Транспорт за месяц','expense','Транспорт');
 await add('400','2026-10-15','Оплата жилья','expense','Жильё');
 for(const [name,amount] of [['Транспорт','120'],['Жильё','450']]){
   await page.goto(base+'/categories');form=page.locator('article.category-row').filter({has:page.getByRole('heading',{name,exact:true})}).locator('form[action="/limit"]');await form.locator('[name=amount]').fill(amount);await submit(form,'Сохранить');
 }
 for(const route of ['overview','operations','categories','payments','analytics']){
  await page.goto(base+'/'+route);await page.screenshot({path:path.join(out,route+'.png'),fullPage:true});
  for(const width of [360,768,1440]){await page.setViewportSize({width,height:1000});assert.ok(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth),`${route}: overflow at ${width}`);}
 }
 await page.setViewportSize({width:360,height:850});await page.goto(base+'/overview');await page.screenshot({path:path.join(out,'mobile.png'),fullPage:true});
 await page.getByRole('button',{name:'Выйти',exact:true}).click();await page.waitForURL('**/login');await page.goto(base+'/overview');assert.match(page.url(),/\/login$/);
 assert.deepEqual(errors,[]);console.log(JSON.stringify({status:'PASS',browser:browser.version(),screenshots:out,viewports:[360,768,1440],checks:['registration','budget','validation','operation edit/delete/cancel','limit','recurring rule','payment replay','link existing operation','delete linked expense','forecasts','logout']},null,2));
 await context.close();
 }finally{await browser.close();}
})().catch(error=>{console.error(error);process.exitCode=1;});
