package fun

import (
"context"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types/events"
"google.golang.org/protobuf/proto"

"cometa-bot/internal/registry"
"cometa-bot/internal/utils"
)

func init() {
registry.Register(&registry.Command{
Name:        "dino",
Aliases:     []string{"dinosaurio", "trex", "chrome"},
Description: "juego del dinosaurio con rola de Vicente Fernández",
Category:    "𝐉⵿𝐮⵿𝐞⵿𝐠⵿𝐨⵿𝐬",
Handler:     Dino,
})
}

const dinoHTML = `<style>
*{box-sizing:border-box;-webkit-tap-highlight-color:transparent;user-select:none;-webkit-user-select:none;margin:0;padding:0}
body{margin:0;background:transparent;font-family:'Segoe UI',Roboto,Arial,sans-serif;color:#eee;overflow:hidden}
.dn-wrap{width:100%;max-width:600px;margin:auto;padding:8px}
.dn-card{background:linear-gradient(135deg,rgba(15,18,28,.98),rgba(30,41,59,.95));border:2px solid rgba(255,107,0,.4);border-radius:20px;overflow:hidden;box-shadow:0 8px 32px rgba(255,107,0,.2)}
.dn-header{padding:10px 14px;background:linear-gradient(90deg,#ff6b00,#ff9500);display:flex;justify-content:space-between;align-items:center;border-bottom:2px solid rgba(255,255,255,.1)}
.dn-brand{display:flex;align-items:center;gap:8px}
.dn-logo{font-size:22px}
.dn-title{font-size:16px;font-weight:900;color:#fff;text-shadow:2px 2px 4px rgba(0,0,0,.5);letter-spacing:1px}
.dn-sub{font-size:9px;color:#fff;font-weight:700;opacity:.9}
.dn-music-btn{background:rgba(255,255,255,.2);border:0;border-radius:50%;width:36px;height:36px;color:#fff;font-size:18px;cursor:pointer;display:flex;align-items:center;justify-content:center}
.dn-music-btn:active{background:rgba(255,255,255,.4)}
.dn-body{padding:12px;text-align:center}
.dn-score{display:flex;justify-content:space-around;margin-bottom:8px;gap:8px}
.dn-box{background:rgba(255,107,0,.15);padding:6px 14px;border-radius:12px;border:1px solid rgba(255,107,0,.3);flex:1}
.dn-box h3{margin:0;font-size:10px;color:#ff9500;letter-spacing:1px;font-weight:700}
.dn-box h1{margin:2px 0 0;font-size:18px;font-family:'Courier New',monospace;color:#fff}
#dnCanvas{background:linear-gradient(180deg,#0b1120 0%,#1a1f3a 100%);border-radius:12px;border:2px solid #ff6b00;touch-action:none;max-width:100%;display:block;margin:0 auto;width:100%;height:auto}
.dn-ctrl{display:flex;gap:8px;justify-content:center;margin-top:10px}
.dn-btn{border:0;border-radius:14px;background:linear-gradient(135deg,#ff6b00,#ff9500);color:#fff;font-size:14px;font-weight:800;cursor:pointer;padding:12px 20px;flex:1;max-width:140px;box-shadow:0 4px 12px rgba(255,107,0,.3)}
.dn-btn:active{transform:scale(.95)}
#dnStatus{margin:8px 0 0;font-size:13px;font-weight:700;color:#ff9500;min-height:18px}
#dnNew{display:none;margin:8px auto 0;padding:10px 22px;border:0;border-radius:12px;background:linear-gradient(135deg,#22c55e,#16a34a);color:#fff;font-size:14px;font-weight:800;cursor:pointer;box-shadow:0 4px 12px rgba(34,197,94,.3)}
.dn-footer{padding:8px;text-align:center;font-size:10px;color:rgba(255,255,255,.5);border-top:1px solid rgba(255,107,0,.2)}
.dn-footer b{color:#ff9500}
.dn-now{font-size:10px;color:#ffd700;margin-top:4px;font-style:italic}
@media(max-width:400px){.dn-title{font-size:14px}.dn-box h1{font-size:16px}.dn-btn{font-size:12px;padding:10px 14px}}
</style>
<div class="dn-wrap">
  <div class="dn-card">
    <div class="dn-header">
      <div class="dn-brand">
        <div class="dn-logo">☄️</div>
        <div><div class="dn-sub">COMETA BOT</div><div class="dn-title">🦖 DINO RUN</div></div>
      </div>
      <button class="dn-music-btn" id="dnMusic">🎵</button>
    </div>
    <div class="dn-body">
      <div class="dn-score">
        <div class="dn-box"><h3>PUNTOS</h3><h1 id="dnScore">0</h1></div>
        <div class="dn-box"><h3>RÉCORD</h3><h1 id="dnBest">0</h1></div>
      </div>
      <canvas id="dnCanvas" width="480" height="160"></canvas>
      <div class="dn-ctrl">
        <button class="dn-btn" id="dnJump">⬆ SALTAR</button>
        <button class="dn-btn" id="dnDuck">⬇ AGACHAR</button>
      </div>
      <p id="dnStatus">Toca SALTAR para empezar 🎮</p>
      <div class="dn-now" id="dnNow">🎺 Rola: El Rey - Vicente Fernández</div>
      <button id="dnNew">🔄 Nuevo juego</button>
    </div>
    <div class="dn-footer">
      ☄️ <b>Cometa Bot</b> • Creado por <b>DevLuis</b> 🇳🇮
    </div>
  </div>
</div>
<script>
(function(){
var cv=document.getElementById('dnCanvas'),ctx=cv.getContext('2d');
var W=cv.width,H=cv.height,GY=130;
var scoreEl=document.getElementById('dnScore'),bestEl=document.getElementById('dnBest');
var statusEl=document.getElementById('dnStatus'),newBtn=document.getElementById('dnNew');
var jumpBtn=document.getElementById('dnJump'),duckBtn=document.getElementById('dnDuck');
var musicBtn=document.getElementById('dnMusic');
var dino={x:40,y:GY-40,w:34,h:40,vy:0,duck:false,onGround:true};
var obs=[],clouds=[{x:120,y:28},{x:300,y:50},{x:430,y:20}];
var speed=220,score=0,best=0,dashOff=0,spawnT=1.0,state='ready',time=0,lastScore=-1;
var audioCtx=null,musicOn=false,musicTimer=null,musicStep=0;

// === NOTA: frecuencias ===
var N={
  A2:110,B2:123.47,Cs3:138.59,D3:146.83,E3:164.81,Fs3:185,A3:220,B3:246.94,
  Cs4:277.18,D4:293.66,E4:329.63,Fs4:369.99,A4:440,B4:493.88,Cs5:554.37,D5:587.33,E5:659.25
};

// === MELODÍA "EL REY" (Vicente Fernández) - compás 6/8 ===
// [nota, duración en corcheas (1=corchea, 2=blanca con puntillo aprox)]
var melodia=[
  // Frase 1: "Yo sé bien que..."
  [N.E4,2],[N.E4,1],[N.Fs4,2],[N.E4,1],
  [N.D4,2],[N.Cs4,1],[N.B3,2],[N.A3,1],
  // Frase 2
  [N.B3,2],[N.Cs4,1],[N.D4,2],[N.E4,1],
  [N.Fs4,3],[N.E4,1],[N.D4,2],
  // Frase 3: "A donde..."
  [N.A4,2],[N.A4,1],[N.B4,2],[N.A4,1],
  [N.Fs4,2],[N.E4,1],[N.D4,2],[N.Cs4,1],
  // Frase 4
  [N.D4,2],[N.E4,1],[N.Fs4,2],[N.E4,1],
  [N.D4,3],[N.Cs4,1],[N.B3,2]
];

// === ACORDES (bajo + rasgueo) ===
var acordes=[
  {root:N.A3,notas:[N.A3,N.Cs4,N.E4]},   // La mayor
  {root:N.D3,notas:[N.D3,N.Fs3,N.A3]},   // Re mayor
  {root:N.E3,notas:[N.E3,N.Gs3,N.B3]},   // Mi mayor (Gs aprox = G#)
  {root:N.A3,notas:[N.A3,N.Cs4,N.E4]}    // La mayor
];

function initAudio(){
  if(audioCtx)return;
  try{audioCtx=new(window.AudioContext||window.webkitAudioContext)()}catch(e){}
}

function playNote(freq,dur,type,vol,delay){
  if(!audioCtx)return;
  delay=delay||0;
  var t=audioCtx.currentTime+delay;
  var osc=audioCtx.createOscillator();
  var gain=audioCtx.createGain();
  osc.type=type;
  osc.frequency.setValueAtTime(freq,t);
  gain.gain.setValueAtTime(0,t);
  gain.gain.linearRampToValueAtTime(vol,t+0.02);
  gain.gain.exponentialRampToValueAtTime(0.001,t+dur);
  osc.connect(gain);
  gain.connect(audioCtx.destination);
  osc.start(t);
  osc.stop(t+dur+0.05);
}

// Trompeta (diente de sierra con filtro)
function playTrompeta(freq,dur){
  if(!audioCtx)return;
  var t=audioCtx.currentTime;
  var osc=audioCtx.createOscillator();
  var osc2=audioCtx.createOscillator();
  var gain=audioCtx.createGain();
  var filter=audioCtx.createBiquadFilter();
  osc.type='sawtooth';
  osc2.type='square';
  osc.frequency.value=freq;
  osc2.frequency.value=freq*1.002;
  filter.type='lowpass';
  filter.frequency.value=1800;
  gain.gain.setValueAtTime(0,t);
  gain.gain.linearRampToValueAtTime(0.12,t+0.05);
  gain.gain.setValueAtTime(0.12,t+dur-0.1);
  gain.gain.linearRampToValueAtTime(0.001,t+dur);
  osc.connect(filter);
  osc2.connect(filter);
  filter.connect(gain);
  gain.connect(audioCtx.destination);
  osc.start(t);osc2.start(t);
  osc.stop(t+dur+0.1);osc2.stop(t+dur+0.1);
}

// Guitarrón (bajo)
function playGuitarron(freq,dur){
  if(!audioCtx)return;
  var t=audioCtx.currentTime;
  var osc=audioCtx.createOscillator();
  var gain=audioCtx.createGain();
  osc.type='triangle';
  osc.frequency.value=freq;
  gain.gain.setValueAtTime(0.25,t);
  gain.gain.exponentialRampToValueAtTime(0.001,t+dur);
  osc.connect(gain);
  gain.connect(audioCtx.destination);
  osc.start(t);osc.stop(t+dur+0.05);
}

// Vihuela (rasgueo)
function playVihuela(nota,dur){
  if(!audioCtx)return;
  var t=audioCtx.currentTime;
  var osc=audioCtx.createOscillator();
  var gain=audioCtx.createGain();
  osc.type='triangle';
  osc.frequency.value=nota*2;
  gain.gain.setValueAtTime(0.08,t);
  gain.gain.exponentialRampToValueAtTime(0.001,t+dur);
  osc.connect(gain);
  gain.connect(audioCtx.destination);
  osc.start(t);osc.stop(t+dur+0.05);
}

// === REPRODUCIR UN COMPÁS DE LA ROLA ===
function playCompas(compasIdx){
  var acorde=acordes[compasIdx%acordes.length];
  var tempo=0.22; // duración de corchea en segundos (~136 BPM en 6/8)
  
  // Guitarrón en tiempos 1 y 4 (patrón ranchero)
  playGuitarron(acorde.root,tempo*1.5);
  setTimeout(function(){playGuitarron(acorde.root*0.5,tempo*1.2)},tempo*3*1000);
  
  // Vihuela: rasgueo en cada corchea
  for(var i=0;i<6;i++){
    setTimeout(function(){
      var n=acorde.notas[i%acorde.notas.length];
      playVihuela(n,tempo*0.8);
    },i*tempo*1000);
  }
  
  // Melodía de trompeta (entra en compases pares)
  if(compasIdx%2===0){
    var frase=melodia.slice((compasIdx*6)%melodia.length,(compasIdx*6+6)%melodia.length+6);
    for(var j=0;j<frase.length&&j<6;j++){
      var nota=frase[j];
      if(nota){
        setTimeout(function(n,d){
          playTrompeta(n[0],d*tempo*0.9);
        },j*tempo*1000,[nota[0],nota[1]]);
      }
    }
  }
}

function startMusica(){
  if(!audioCtx||!musicOn)return;
  if(musicTimer)clearInterval(musicTimer);
  var tempo=0.22;
  var compasDur=6*tempo*1000; // duración de un compás en ms
  musicStep=0;
  playCompas(0);
  musicTimer=setInterval(function(){
    if(!musicOn||state!=='run'){
      clearInterval(musicTimer);
      return;
    }
    musicStep=(musicStep+1)%4;
    playCompas(musicStep);
  },compasDur);
}

function toggleMusic(){
  initAudio();
  musicOn=!musicOn;
  musicBtn.textContent=musicOn?'🔊':'🎵';
  if(musicOn){
    if(state==='run')startMusica();
  }else{
    if(musicTimer)clearInterval(musicTimer);
  }
}

function playJump(){
  if(!audioCtx)return;
  playNote(400,0.1,'square',0.08);
  setTimeout(function(){playNote(600,0.1,'square',0.06)},50);
}

function playScore(){
  if(!audioCtx)return;
  playNote(880,0.08,'sine',0.05);
  setTimeout(function(){playNote(1100,0.1,'sine',0.04)},60);
}

function playGameOver(){
  if(!audioCtx)return;
  // Acorde dramático de trompeta
  playTrompeta(N.A3,0.4);
  setTimeout(function(){playTrompeta(N.E3,0.5)},150);
  setTimeout(function(){playTrompeta(N.A2,0.6)},300);
}

function reset(){
  obs.length=0;speed=220;score=0;dashOff=0;spawnT=1.0;time=0;lastScore=-1;
  dino.y=GY-40;dino.h=40;dino.w=34;dino.vy=0;dino.duck=false;dino.onGround=true;
  state='run';scoreEl.textContent='0';statusEl.textContent='¡Corre! 🏃';newBtn.style.display='none';
  if(musicOn)startMusica();
}

function jump(){
  initAudio();
  if(state==='over'){reset();return}
  if(state==='ready'){state='run';statusEl.textContent='¡Corre! 🏃';if(musicOn)startMusica()}
  if(dino.onGround){dino.vy=-620;dino.onGround=false;playJump()}
}

function setDuck(v){dino.duck=v;if(v&&!dino.onGround){dino.vy+=320}}

function spawn(){
  var r=Math.random(),o;
  if(score>150&&r<0.22){o={t:'b',x:W+20,y:(Math.random()<0.5?76:104),w:26,h:16}}
  else if(r<0.5){o={t:'c1',x:W+20,y:GY-30,w:14,h:30}}
  else if(r<0.8){o={t:'c2',x:W+20,y:GY-40,w:18,h:40}}
  else{o={t:'c3',x:W+20,y:GY-34,w:32,h:34}}
  obs.push(o);
}

function update(dt){
  time+=dt;speed=Math.min(520,speed+8*dt);dashOff=(dashOff+speed*dt)%24;
  for(var i=0;i<clouds.length;i++){var c=clouds[i];c.x-=speed*0.25*dt;if(c.x<-40){c.x=W+40;c.y=16+Math.random()*45}}
  spawnT-=dt;if(spawnT<=0){spawn();spawnT=0.55+Math.random()*0.9*(260/speed)+0.35}
  dino.vy+=2000*dt;dino.y+=dino.vy*dt;
  var sh=dino.duck?24:40;
  if(dino.y>=GY-sh){dino.y=GY-sh;dino.vy=0;dino.onGround=true}else{dino.onGround=false}
  dino.h=sh;dino.w=dino.duck?44:34;
  var wI=0;for(var j=0;j<obs.length;j++){var o=obs[j];o.x-=speed*dt;if(o.x+o.w>-10){obs[wI++]=o}}obs.length=wI;
  var dx=dino.x+4,dy=dino.y+4,dw=dino.w-8,dh=dino.h-8;
  for(var k=0;k<obs.length;k++){var ob=obs[k];var ox=ob.x+3,oy=ob.y+3,ow=ob.w-6,oh=ob.h-6;if(dx<ox+ow&&dx+dw>ox&&dy<oy+oh&&dy+dh>oy){gameOver();break}}
  score+=speed*dt*0.05;var sc=Math.floor(score);if(sc!==lastScore){lastScore=sc;scoreEl.textContent=sc;if(sc%50===0&&sc>0)playScore()}
}

function gameOver(){
  state='over';if(Math.floor(score)>best){best=Math.floor(score);bestEl.textContent=best}
  statusEl.textContent='💀 Game Over • Puntos: '+Math.floor(score)+' • Toca para reiniciar';
  newBtn.style.display='block';playGameOver();
  if(musicTimer)clearInterval(musicTimer);
}

function drawDino(){
  ctx.fillStyle='#ff6b00';var x=dino.x,y=dino.y;
  if(dino.duck){
    ctx.fillRect(x,y+8,34,16);ctx.fillRect(x+30,y,14,12);
    ctx.fillStyle='#0b1120';ctx.fillRect(x+38,y+3,3,3);
    ctx.fillStyle='#ff6b00';var lp=Math.floor(time*12)%2;
    if(lp){ctx.fillRect(x+4,y+24,6,0);ctx.fillRect(x+22,y+24,6,0)}
  }else{
    ctx.fillRect(x,y+10,26,22);ctx.fillRect(x+12,y,22,16);ctx.fillRect(x-6,y+14,6,8);
    ctx.fillStyle='#0b1120';ctx.fillRect(x+26,y+4,3,3);
    ctx.fillStyle='#ff6b00';var lp2=Math.floor(time*12)%2;
    if(dino.onGround){
      if(lp2){ctx.fillRect(x+4,y+32,6,8);ctx.fillRect(x+16,y+32,6,5)}
      else{ctx.fillRect(x+4,y+32,6,5);ctx.fillRect(x+16,y+32,6,8)}
    }else{ctx.fillRect(x+4,y+32,6,6);ctx.fillRect(x+16,y+32,6,6)}
  }
}

function draw(){
  ctx.fillStyle='#0b1120';ctx.fillRect(0,0,W,H);
  ctx.fillStyle='#20304a';for(var i=0;i<clouds.length;i++){var c=clouds[i];ctx.fillRect(c.x,c.y,34,8);ctx.fillRect(c.x+6,c.y-6,20,6)}
  ctx.fillStyle='#ff6b00';ctx.fillRect(0,GY,W,2);
  ctx.fillStyle='#243447';for(var x=-dashOff;x<W;x+=24){ctx.fillRect(x,GY+6,12,2)}
  for(var j=0;j<obs.length;j++){
    var o=obs[j];
    if(o.t==='b'){
      ctx.fillStyle='#ff9f43';ctx.fillRect(o.x,o.y+5,o.w-6,8);
      var up=Math.sin(time*18)>0;ctx.fillRect(o.x+8,up?o.y:o.y+10,9,5);
      ctx.fillStyle='#ffd166';ctx.fillRect(o.x+o.w-7,o.y+7,7,3);
    }else{
      ctx.fillStyle='#22c55e';
      if(o.t==='c1'){ctx.fillRect(o.x+4,o.y,6,o.h);ctx.fillRect(o.x,o.y+8,4,8);ctx.fillRect(o.x+10,o.y+12,4,8)}
      else if(o.t==='c2'){ctx.fillRect(o.x+6,o.y,6,o.h);ctx.fillRect(o.x,o.y+10,5,10);ctx.fillRect(o.x+13,o.y+14,5,10)}
      else{ctx.fillRect(o.x+2,o.y+4,5,o.h-4);ctx.fillRect(o.x+13,o.y,6,o.h);ctx.fillRect(o.x+24,o.y+6,5,o.h-6)}
    }
  }
  drawDino();
}

var last=performance.now();
function frame(now){var dt=(now-last)/1000;last=now;if(dt>0.05){dt=0.05}if(state==='run'){update(dt)}draw();requestAnimationFrame(frame)}
requestAnimationFrame(frame);
document.addEventListener('visibilitychange',function(){last=performance.now()});

cv.addEventListener('touchstart',function(e){e.preventDefault();jump()},{passive:false});
cv.addEventListener('mousedown',function(){jump()});
jumpBtn.addEventListener('touchstart',function(e){e.preventDefault();jump()},{passive:false});
jumpBtn.addEventListener('click',function(){jump()});
duckBtn.addEventListener('touchstart',function(e){e.preventDefault();setDuck(true)},{passive:false});
duckBtn.addEventListener('touchend',function(){setDuck(false)});
duckBtn.addEventListener('mousedown',function(){setDuck(true)});
duckBtn.addEventListener('mouseup',function(){setDuck(false)});
duckBtn.addEventListener('mouseleave',function(){setDuck(false)});
newBtn.addEventListener('click',function(){reset()});
musicBtn.addEventListener('click',function(){toggleMusic()});

document.addEventListener('keydown',function(e){
  if(e.key==='ArrowUp'||e.key===' '){e.preventDefault();jump()}
  if(e.key==='ArrowDown'){setDuck(true)}
});
document.addEventListener('keyup',function(e){if(e.key==='ArrowDown'){setDuck(false)}});
})();
</script>`

func Dino(evt *events.Message, client *whatsmeow.Client, _ []string) {
msg := utils.BuildBotForwardedMessage(dinoHTML)
if msg == nil {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("❌ Esta versión de WhatsApp no soporta juegos HTML interactivos."),
ContextInfo: &waE2E.ContextInfo{
QuotedMessage: evt.Message,
},
},
})
return
}

client.SendMessage(context.Background(), evt.Info.Chat, msg)
}
