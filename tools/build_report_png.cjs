// Rasterize editable report SVGs without browser or HTML templates.
const fs = require('fs');
const path = require('path');
const sharp = require('sharp');
const dir = path.resolve(__dirname, '../reports/lab3/figures');
(async () => {
  const names = fs.readdirSync(dir).filter(n => n.endsWith('.svg')).sort();
  for (const name of names) {
    await sharp(path.join(dir, name)).png().toFile(path.join(dir, name.replace(/\.svg$/, '.png')));
  }
  console.log(`Rendered ${names.length} PNG illustrations`);
})().catch(e => { console.error(e); process.exit(1); });
