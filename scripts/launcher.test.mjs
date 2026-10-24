import {test} from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import vm from 'node:vm';
const source=fs.readFileSync(new URL('../web/launcher.js',import.meta.url),'utf8');
function boot(){
 const events=new Map(),messages=[],timers=new Map();let sequence=0,focused=0,exits=0;
 const element=(dataset={})=>({dataset,style:{},events:new Map(),addEventListener(type,fn){this.events.set(type,fn);},setPointerCapture(){},focus(){focused++;}});
 const frame=element({src:'game.html'}),status=element(),canvas=element();frame.contentDocument={querySelector:()=>canvas};frame.contentWindow={postMessage:m=>messages.push(m)};frame.requestFullscreen=async()=>{};
 const nodes={game:frame,status,restart:element(),focus:element(),fullscreen:element()};
 const buttons=['left','right','up','down','attack','use'].map(control=>element({control}));
 const document={hidden:false,fullscreenElement:null,getElementById:id=>nodes[id],querySelectorAll:()=>buttons,addEventListener:(type,fn)=>events.set('document-'+type,fn),exitFullscreen:async()=>{exits++;}};
 const context={URL,crypto:{randomUUID:()=>String(++sequence)},location:{href:'https://example.test/',origin:'https://example.test'},document,setTimeout:fn=>{const id=timers.size+1;timers.set(id,fn);return id;},clearTimeout:id=>timers.delete(id),addEventListener:(type,fn)=>events.set(type,fn)};context.window=context;
 vm.runInNewContext(source,context);
 const session=()=>new URL(frame.src).searchParams.get('session');
 const send=(data,overrides={})=>events.get('message')({origin:context.location.origin,source:frame.contentWindow,data:{session:session(),...data},...overrides});
 return {context,nodes,frame,status,buttons,events,messages,timers,session,send,focused:()=>focused,exits:()=>exits};
}
test('autostart trusts only current frame/origin/session and focuses after Go readiness',()=>{
 const r=boot();assert.equal(r.status.dataset.state,'loading');
 for(const override of [{origin:'https://foreign.test'},{source:{} }])r.send({type:'app-ready'},override);
 r.send({type:'app-ready',session:'stale'});assert.equal(r.focused(),0);
 r.send({type:'app-ready',width:1920,height:1440});assert.equal(r.status.dataset.state,'ready');assert.equal(r.frame.style.aspectRatio,'1920/1440');assert.equal(r.focused(),2);
});
test('slow loading permits late success; errors remain sticky until a fresh restart',()=>{
 const r=boot();[...r.timers.values()][0]();assert.match(r.status.textContent,/Still loading/);r.send({type:'app-ready'});assert.equal(r.status.dataset.state,'ready');
 r.send({type:'app-error',message:'bad resource'});assert.match(r.status.textContent,/bad resource/);r.send({type:'app-ready'});r.send({type:'app-error',message:'secondary'});assert.match(r.status.textContent,/bad resource/);
 const old=r.session();r.nodes.restart.events.get('click')();assert.notEqual(r.session(),old);assert.equal(r.status.dataset.state,'loading');r.send({type:'app-error',session:old,message:'late error'});assert.equal(r.status.dataset.state,'loading');
});
test('pointer/keyboard controls release on cancellation, blur and hidden pages',()=>{
 const r=boot();r.send({type:'app-ready'});const button=r.buttons[1],event={button:0,pointerId:4,key:' ',repeat:false,preventDefault(){}};
 button.events.get('pointerdown')(event);assert.equal(r.messages.at(-1).key,'right');assert.equal(r.messages.at(-1).down,true);
 for(const type of ['pointerup','pointercancel','lostpointercapture','blur']){button.events.get(type)(event);assert.equal(r.messages.at(-1).down,false);}
 button.events.get('keydown')(event);assert.equal(r.messages.at(-1).down,true);button.events.get('keyup')(event);assert.equal(r.messages.at(-1).down,false);
 r.events.get('blur')();assert.ok(r.messages.slice(-6).every(m=>m.down===false));r.context.document.hidden=true;r.events.get('document-visibilitychange')();assert.ok(r.messages.slice(-6).every(m=>m.down===false));
});
test('fullscreen denial is recoverable and only the trusted frame can request exit',async()=>{
 const r=boot();r.send({type:'app-ready'});r.frame.requestFullscreen=async()=>{throw new Error('denied')};await r.nodes.fullscreen.events.get('click')();assert.match(r.status.textContent,/unavailable/);assert.equal(r.status.dataset.state,'ready');
 r.context.document.fullscreenElement=r.frame;r.send({type:'app-exit-fullscreen'},{origin:'https://foreign.test'});assert.equal(r.exits(),0);r.send({type:'app-exit-fullscreen'});assert.equal(r.exits(),1);
 r.events.get('document-fullscreenchange')();assert.equal(r.messages.at(-1).type,'app-display');assert.equal(r.messages.at(-1).fullscreen,true);
});
