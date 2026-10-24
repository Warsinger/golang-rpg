const frame = document.getElementById('game');
const status = document.getElementById('status');
const keys = ['left','right','up','down','attack','use'];
let session = '', ready = false, timeout;
function focusGame() {
  frame.focus({preventScroll:true});
  const canvas = frame.contentDocument?.querySelector('canvas');
  if (canvas) {canvas.tabIndex=0;canvas.focus({preventScroll:true});}
}
function sendInput(key,down) {if (ready) frame.contentWindow.postMessage({type:'app-input',session,key,down},location.origin);}
function releaseInput() {for (const key of keys) sendInput(key,false);}
function launch() {
  releaseInput();ready=false;session=crypto.randomUUID();clearTimeout(timeout);
  status.dataset.state='loading';status.textContent='Loading… The first visit downloads the game.';
  const url=new URL(frame.dataset.src,location.href);url.searchParams.set('session',session);frame.src=url.href;
  timeout=setTimeout(()=>{if(!ready && status.dataset.state==='loading') status.textContent='Still loading. If the game does not appear, use Restart scene to retry.';},60000);
}
window.addEventListener('message',event=>{
  if(event.origin!==location.origin || event.source!==frame.contentWindow || event.data?.session!==session) return;
  if(event.data.type==='app-exit-fullscreen'){if(document.fullscreenElement===frame)document.exitFullscreen().catch(()=>{});return;}
  if(status.dataset.state==='runtime-error') return;
  if(event.data.type==='app-ready') {
    ready=true;clearTimeout(timeout);status.dataset.state='ready';status.textContent='Ready. Use the keyboard or the controls below.';
    const {width,height}=event.data;
    if(Number.isInteger(width)&&Number.isInteger(height)&&width>0&&height>0&&width<=8192&&height<=8192)frame.style.aspectRatio=width+'/'+height;
    frame.contentWindow.postMessage({type:'app-display',session,fullscreen:document.fullscreenElement===frame},location.origin);
    focusGame();
  } else if(event.data.type==='app-error') {
    releaseInput();ready=false;clearTimeout(timeout);status.dataset.state='runtime-error';status.textContent='Could not run. Restart scene to retry. '+String(event.data.message||'Unknown error');
  }
});
document.getElementById('restart').addEventListener('click',launch);
document.getElementById('focus').addEventListener('click',focusGame);
document.getElementById('fullscreen').addEventListener('click',async()=>{
  try {await frame.requestFullscreen();focusGame();}
  catch {if(status.dataset.state!=='runtime-error')status.textContent='Fullscreen is unavailable here. You can still play in the page.';}
});
for(const button of document.querySelectorAll('[data-control]')) {
  button.addEventListener('pointerdown',event=>{if(event.button!==0)return;event.preventDefault();button.setPointerCapture(event.pointerId);sendInput(button.dataset.control,true);});
  for(const type of ['pointerup','pointercancel','lostpointercapture']) button.addEventListener(type,()=>sendInput(button.dataset.control,false));
  button.addEventListener('keydown',event=>{if(['Enter',' '].includes(event.key)){event.preventDefault();if(!event.repeat)sendInput(button.dataset.control,true);}});
  button.addEventListener('keyup',event=>{if(['Enter',' '].includes(event.key)){event.preventDefault();sendInput(button.dataset.control,false);}});
  button.addEventListener('blur',()=>sendInput(button.dataset.control,false));
}
window.addEventListener('blur',releaseInput);
document.addEventListener('visibilitychange',()=>{if(document.hidden)releaseInput();});
launch();

document.addEventListener('fullscreenchange',()=>{frame.contentWindow.postMessage({type:'app-display',session,fullscreen:document.fullscreenElement===frame},location.origin);if(ready)focusGame();});
