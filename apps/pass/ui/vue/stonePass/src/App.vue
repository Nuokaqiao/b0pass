<script setup></script>

<template>
  <div class="app-root">
    <div class="ambient" aria-hidden="true">
      <span class="wash"></span>
      <span class="blob blob-a"></span>
      <span class="blob blob-b"></span>
      <span class="blob blob-c"></span>
      <span class="blob blob-d"></span>
      <span class="blob blob-e"></span>
      <span class="band band-a"></span>
      <span class="band band-b"></span>
      <span class="ring ring-a"></span>
      <span class="ring ring-b"></span>
      <span class="ring ring-c"></span>
      <span class="shape shape-a"></span>
      <span class="shape shape-b"></span>
      <span class="shape shape-c"></span>
      <span class="grid"></span>
      <span class="dots"></span>
      <span class="wave"></span>
    </div>
    <div class="app-content">
      <router-view />
    </div>
  </div>
</template>

<style scoped>
.app-root {
  position: relative;
  min-height: 100vh;
  overflow-x: hidden;
}

.app-content {
  position: relative;
  z-index: 1;
}

.ambient {
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  overflow: hidden;
  background: linear-gradient(160deg, #d7e4ee 0%, #c9dbe8 38%, #bfd4c8 72%, #c5d6e6 100%);
}

.wash {
  position: absolute;
  inset: 0;
  background:
    radial-gradient(ellipse 90% 70% at 15% 10%, rgba(90, 168, 190, 0.45), transparent 55%),
    radial-gradient(ellipse 80% 60% at 90% 20%, rgba(95, 168, 130, 0.38), transparent 50%),
    radial-gradient(ellipse 70% 55% at 50% 95%, rgba(110, 148, 190, 0.4), transparent 55%),
    radial-gradient(ellipse 50% 40% at 70% 55%, rgba(150, 180, 200, 0.35), transparent 60%);
}

.blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(40px);
  opacity: 0.75;
  animation: drift 16s ease-in-out infinite;
}

.blob-a {
  width: 520px;
  height: 520px;
  left: -160px;
  top: -100px;
  background: radial-gradient(circle, #8ec8d8 0%, rgba(142, 200, 216, 0) 68%);
}

.blob-b {
  width: 460px;
  height: 460px;
  right: -140px;
  top: 5%;
  background: radial-gradient(circle, #8fbfa4 0%, rgba(143, 191, 164, 0) 68%);
  animation-delay: -5s;
}

.blob-c {
  width: 540px;
  height: 540px;
  left: 20%;
  bottom: -220px;
  background: radial-gradient(circle, #9bb4d0 0%, rgba(155, 180, 208, 0) 70%);
  animation-delay: -10s;
}

.blob-d {
  width: 280px;
  height: 280px;
  left: 42%;
  top: 28%;
  background: radial-gradient(circle, #a8cfe0 0%, rgba(168, 207, 224, 0) 70%);
  opacity: 0.55;
  animation-delay: -3s;
}

.blob-e {
  width: 320px;
  height: 320px;
  right: 18%;
  bottom: 22%;
  background: radial-gradient(circle, #a3c9b4 0%, rgba(163, 201, 180, 0) 70%);
  opacity: 0.5;
  animation-delay: -8s;
}

.band {
  position: absolute;
  left: -10%;
  width: 120%;
  height: 180px;
  opacity: 0.22;
  filter: blur(2px);
  transform: rotate(-8deg);
}

.band-a {
  top: 22%;
  background: linear-gradient(90deg, transparent, rgba(15, 110, 140, 0.35), rgba(45, 106, 79, 0.28), transparent);
  animation: sway 20s ease-in-out infinite;
}

.band-b {
  top: 68%;
  height: 140px;
  transform: rotate(6deg);
  background: linear-gradient(90deg, transparent, rgba(70, 120, 160, 0.3), rgba(15, 110, 140, 0.25), transparent);
  animation: sway 24s ease-in-out infinite reverse;
}

.ring {
  position: absolute;
  border-radius: 50%;
  border: 1.5px solid rgba(15, 110, 140, 0.22);
  animation: pulse 12s ease-in-out infinite;
}

.ring-a {
  width: 320px;
  height: 320px;
  right: 6%;
  bottom: 12%;
  border-color: rgba(45, 106, 79, 0.26);
}

.ring-b {
  width: 180px;
  height: 180px;
  left: 8%;
  top: 38%;
  animation-delay: -4s;
}

.ring-c {
  width: 120px;
  height: 120px;
  right: 28%;
  top: 18%;
  border-color: rgba(90, 140, 170, 0.28);
  animation-delay: -7s;
}

.shape {
  position: absolute;
  border: 1.5px solid rgba(15, 110, 140, 0.18);
  background: rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(2px);
}

.shape-a {
  width: 90px;
  height: 90px;
  left: 12%;
  bottom: 16%;
  border-radius: 18px;
  transform: rotate(18deg);
  animation: floaty 14s ease-in-out infinite;
}

.shape-b {
  width: 64px;
  height: 64px;
  right: 14%;
  top: 36%;
  border-radius: 50%;
  border-style: dashed;
  animation: floaty 11s ease-in-out infinite reverse;
}

.shape-c {
  width: 110px;
  height: 48px;
  left: 55%;
  top: 12%;
  border-radius: 999px;
  transform: rotate(-12deg);
  animation: floaty 17s ease-in-out infinite;
  animation-delay: -2s;
}

.grid {
  position: absolute;
  inset: 0;
  opacity: 0.55;
  background-image:
    linear-gradient(rgba(15, 110, 140, 0.1) 1px, transparent 1px),
    linear-gradient(90deg, rgba(15, 110, 140, 0.1) 1px, transparent 1px);
  background-size: 36px 36px;
}

.dots {
  position: absolute;
  inset: 0;
  opacity: 0.55;
  background-image: radial-gradient(rgba(15, 90, 120, 0.22) 1.2px, transparent 1.2px);
  background-size: 18px 18px;
}

.wave {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  height: 42%;
  opacity: 0.35;
  background:
    repeating-linear-gradient(
      -12deg,
      transparent 0 18px,
      rgba(255, 255, 255, 0.12) 18px 20px
    ),
    linear-gradient(180deg, transparent, rgba(90, 130, 160, 0.18));
  mask-image: linear-gradient(180deg, transparent, #000 40%);
}

@keyframes drift {
  0%,
  100% {
    transform: translate(0, 0) scale(1);
  }
  50% {
    transform: translate(20px, -16px) scale(1.05);
  }
}

@keyframes pulse {
  0%,
  100% {
    transform: scale(1);
    opacity: 0.65;
  }
  50% {
    transform: scale(1.07);
    opacity: 1;
  }
}

@keyframes sway {
  0%,
  100% {
    transform: rotate(-8deg) translateX(0);
  }
  50% {
    transform: rotate(-6deg) translateX(3%);
  }
}

@keyframes floaty {
  0%,
  100% {
    transform: translateY(0) rotate(18deg);
  }
  50% {
    transform: translateY(-14px) rotate(22deg);
  }
}

.shape-b {
  animation-name: floaty-round;
}

.shape-c {
  animation-name: floaty-pill;
}

@keyframes floaty-round {
  0%,
  100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-12px);
  }
}

@keyframes floaty-pill {
  0%,
  100% {
    transform: translateY(0) rotate(-12deg);
  }
  50% {
    transform: translateY(-10px) rotate(-8deg);
  }
}

@media (prefers-reduced-motion: reduce) {
  .blob,
  .ring,
  .band,
  .shape {
    animation: none;
  }
}
</style>
