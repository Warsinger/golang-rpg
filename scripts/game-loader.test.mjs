import {test} from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import vm from 'node:vm';
import {gzipSync} from 'node:zlib';
const source=fs.readFileSync(new URL('../web/game.js',import.meta.url),'utf8');
const wasm=Uint8Array.from([0,97,115,109,1,0,0,0]);
async function run({bytes=wasm,ok=true,decompress=true,stop=false}={}) {
 const events=new Map(),messages=[];let instantiated;
 const parent={postMessage:message=>messages.push(message)};
 const context={URLSearchParams,Uint8Array,Response,Blob,WebAssembly:{instantiate:async data=>{instantiated=new Uint8Array(data);return WebAssembly.instantiate(data,{})}},DecompressionStream:decompress?DecompressionStream:undefined,location:{origin:'https://example.test',search:'?session=test'},parent,document:{getElementById:()=>({hidden:true,addEventListener:()=>{}}),hidden:false,addEventListener:(type,fn)=>events.set('document-'+type,fn)},addEventListener:(type,fn)=>events.set(type,fn),fetch:async()=>new Response(bytes,{status:ok?200:404})};
 context.window=context;
 context.Go=class{importObject={};run(){context.appReady(1920,1440);return stop?Promise.reject(new Error('stopped')):new Promise(()=>{});}};
 vm.runInNewContext(source,context);
 for(let i=0;i<100&&!messages.some(m=>m.type===(stop||!ok||!decompress?'app-error':'app-ready'));i++) await new Promise(resolve=>setTimeout(resolve,5));
 return {context,events,messages,instantiated,parent};
}
test('loads raw Wasm and reports the session only when Go is ready',async()=>{
 const r=await run();assert.deepEqual(r.instantiated,wasm);assert.equal(r.messages[0].type,'app-ready');assert.equal(r.messages[0].session,'test');
});
test('loads gzip Wasm and handles browser-decoded responses too',async()=>{
 const gzip=await run({bytes:gzipSync(wasm)});assert.deepEqual(gzip.instantiated,wasm);assert.equal(gzip.messages[0].type,'app-ready');
});
test('download, missing decompression and runtime failures report useful errors',async()=>{
 for(const [options,pattern] of [[{ok:false},/HTTP 404/],[{bytes:gzipSync(wasm),decompress:false},/cannot decompress/],[{stop:true},/stopped/]]) {
  const r=await run(options);assert.match(r.messages.find(m=>m.type==='app-error')?.message,pattern);
 }
});
test('input bridge rejects foreign/stale messages, latches pickup, and releases on blur',async()=>{
 const r=await run();const message=r.events.get('message');
 const send=(data,overrides={})=>message({origin:'https://example.test',source:r.parent,data:{type:'app-input',session:'test',...data},...overrides});
 send({key:'left',down:true},{source:{}});assert.equal(r.context.appInput.left,false);
 send({key:'left',down:true},{origin:'https://foreign.test'});assert.equal(r.context.appInput.left,false);
 send({key:'left',down:true,session:'stale'});assert.equal(r.context.appInput.left,false);
 send({key:'left',down:true});assert.equal(r.context.appInput.left,true);
 send({key:'use',down:true});send({key:'use',down:false});assert.equal(r.context.appInput.use,true);
 for(const key of ['up','down','attack']) {send({key,down:true});assert.equal(r.context.appInput[key],true);send({key,down:false});assert.equal(r.context.appInput[key],false);}
 send({key:'right',down:'yes'});assert.equal(r.context.appInput.right,false);
 send({key:'__proto__',down:true});
 r.events.get('blur')();assert.deepEqual({...r.context.appInput},{left:false,right:false,up:false,down:false,attack:false,use:false});
});
