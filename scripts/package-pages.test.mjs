import {test} from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import {gunzipSync} from 'node:zlib';
import {packagePages} from './package-pages.mjs';
function fixture(t) {
 const root=fs.mkdtempSync(path.join(os.tmpdir(),'game-pages-test-'));
 t.after(()=>fs.rmSync(root,{recursive:true,force:true}));
 const source=path.join(root,'source');fs.mkdirSync(source);
 for(const file of ['index.html','game.html','launcher.js','game.js','styles.css']) fs.copyFileSync(new URL('../web/'+file,import.meta.url),path.join(source,file));
 fs.writeFileSync(path.join(source,'wasm_exec.js'),'// matching Go runtime');
 fs.writeFileSync(path.join(source,'game.wasm'),Buffer.from([0,97,115,109,1,0,0,0]));
 return {root,source};
}
test('packages actual HTML/module references, gzip bytes, revision and cache policy',t=>{
 const {root,source}=fixture(t);const output=path.join(root,'pages');
 const hash=packagePages(source,output,'abcdef');const assets=path.join(output,'assets',hash);
 assert.deepEqual(gunzipSync(fs.readFileSync(path.join(assets,'game.wasm.gz'))),fs.readFileSync(path.join(source,'game.wasm')));
 assert.match(fs.readFileSync(path.join(assets,'game.js'),'utf8'),/fetch\("game.wasm.gz"\)/);
 const index=fs.readFileSync(path.join(output,'index.html'),'utf8');
 for(const file of ['styles.css','launcher.js','game.html']) assert.ok(index.includes('assets/'+hash+'/'+file));
 for(const match of fs.readFileSync(path.join(assets,'game.html'),'utf8').matchAll(/src="([^"]+)"/g)) assert.ok(fs.existsSync(path.join(assets,match[1])));
 assert.deepEqual(JSON.parse(fs.readFileSync(path.join(output,'version.json'))),{revision:'abcdef',assets:'assets/'+hash});
 const headers=fs.readFileSync(path.join(output,'_headers'),'utf8');assert.match(headers,/no-cache/);assert.match(headers,/immutable/);assert.match(headers,/application\/gzip/);
 assert.ok(fs.existsSync(path.join(output,'404.html')));
 assert.throws(()=>packagePages(source,output),/EEXIST/);
});
test('asset paths change when implementation changes and remain deterministic',t=>{
 const {root,source}=fixture(t);
 const a=packagePages(source,path.join(root,'a'),'one');const b=packagePages(source,path.join(root,'b'),'two');assert.equal(a,b);
 fs.appendFileSync(path.join(source,'launcher.js'),'\n// changed');
 assert.notEqual(packagePages(source,path.join(root,'c')),a);
});
test('refuses missing references and overlarge assets before creating output',t=>{
 const {root,source}=fixture(t);fs.writeFileSync(path.join(source,'game.js'),'bad loader');
 assert.throws(()=>packagePages(source,path.join(root,'bad')),/fetch reference/);assert.equal(fs.existsSync(path.join(root,'bad')),false);
 fs.copyFileSync(new URL('../web/game.js',import.meta.url),path.join(source,'game.js'));
 fs.writeFileSync(path.join(source,'wasm_exec.js'),Buffer.alloc(25*1024*1024+1));
 assert.throws(()=>packagePages(source,path.join(root,'large')),/25 MiB/);
});

