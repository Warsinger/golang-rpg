"use strict";
const session=new URLSearchParams(location.search).get('session');
const keys=['left','right','up','down','attack','use'];
const input=window.appInput=Object.fromEntries(keys.map(key=>[key,false]));
const report=(type,message,extra={})=>parent.postMessage({type,session,message,...extra},location.origin);
window.appReady=(width,height)=>{
  report('app-ready',undefined,{width,height});
  if(parent===window){const canvas=document.querySelector('canvas');if(canvas){canvas.tabIndex=0;canvas.focus({preventScroll:true});}}
};
window.appError=message=>report('app-error',String(message));
window.addEventListener('error',event=>window.appError(event.error?.message||event.message));
window.addEventListener('unhandledrejection',event=>window.appError(event.reason?.message||event.reason));
const exitFullscreen=document.getElementById('exit-fullscreen');
exitFullscreen.addEventListener('click',()=>report('app-exit-fullscreen'));
window.addEventListener('message',event=>{
  if(event.origin!==location.origin||event.source!==parent||event.data?.session!==session)return;
  if(event.data.type==='app-display'){exitFullscreen.hidden=event.data.fullscreen!==true;return;}
  if(event.data.type!=='app-input')return;
  const {key,down}=event.data;
  if(!keys.includes(key)||typeof down!=='boolean')return;
  if(key!=='use'||down)input[key]=down;
});
function release(){for(const key of keys)input[key]=false;}
window.addEventListener('blur',release);
document.addEventListener('visibilitychange',()=>{if(document.hidden)release();});
window.addEventListener('keydown',event=>{if(event.code==='Escape')report('app-exit-fullscreen');if(['Space','ArrowUp','ArrowDown','ArrowLeft','ArrowRight'].includes(event.code))event.preventDefault();});
(async()=>{
  try {
    const go=new Go();
    const response=await fetch("game.wasm");
    if(!response.ok)throw new Error('Game download failed: HTTP '+response.status);
    let bytes=await response.arrayBuffer();
    const magic=new Uint8Array(bytes,0,Math.min(2,bytes.byteLength));
    if(magic[0]===0x1f&&magic[1]===0x8b){
      if(typeof DecompressionStream==='undefined')throw new Error('This browser cannot decompress the game. Try a current browser.');
      bytes=await new Response(new Blob([bytes]).stream().pipeThrough(new DecompressionStream('gzip'))).arrayBuffer();
    }
    const result=await WebAssembly.instantiate(bytes,go.importObject);
    await go.run(result.instance);
    throw new Error('Game stopped. Restart scene to run again.');
  } catch(error){window.appError(error.message||error);}
})();
