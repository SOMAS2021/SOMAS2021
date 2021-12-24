import React from 'react';
import '../App.css';
import '../style/Hero.css';
import Video from '../videos/platform.mov';

function Hero() {
  return (
    <div className="hero-container">
      <video src={Video} autoPlay loop muted />
      <h1>SELF ORGANISING MULTI AGENT SYSTEMS 2021</h1>
      <p>The Platform-Playing Platform</p>
    </div>
  );
}

export default Hero;
