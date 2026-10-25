// Syntax check with the locally installed Ellidiss 1.0.0 parser.
// Usage: node tools/check_sysml.cjs /path/to/ellidiss.../sysml-extension
const fs=require('fs'),path=require('path');
const base=process.argv[2];
if(!base){console.error('Pass the directory containing the Ellidiss parser.js and sysml.js');process.exit(2);}
const pa=require(path.resolve(base,'parser.js')).$_$;
const sy=require(path.resolve(base,'sysml.js')).$_$;
const filename=path.resolve(__dirname,'../srs/sysml/balance.sysml');
const text=fs.readFileSync(filename,'utf8'),input=new pa.p(filename,text),parser=sy.a(input);
let tree;try{tree=sy.b(parser);}catch(e){console.error(String(e));process.exit(1);}
const errors=parser.ku_1;
if(!tree||input.gu_1<text.trimEnd().length||errors.length){console.error(JSON.stringify({parsed:!!tree,consumed:input.gu_1,length:text.length,errors}));process.exit(1);}
console.log('PASS: SysML v2 syntax, entire model parsed. Semantic validation is not performed.');
