const fs=require('fs'),path=require('path'),sharp=require('sharp');
(async()=>{const root=path.resolve(__dirname,'..');for(const f of fs.readdirSync(root+'/srs/sysml/views').filter(f=>f.endsWith('.svg')))await sharp(root+'/srs/sysml/views/'+f,{density:144}).png().toFile(root+'/srs/img/'+f.replace('.svg','.png'));console.log('5 PNG projections rendered');})();
