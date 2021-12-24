import React from 'react';
import '../../style/Footer.css';

function Footer() {
  return (
    <div className="footer-container">
      <h1>Imperial College London Electrical & Electronic Engineering Department</h1>
      <p> Self Organising Multi Agent Systems</p>
      <div className="social-icons">
        <a href="https://github.com/SOMAS2021" className="icon-holder" target="_blank" rel="noreferrer">
          Github <i class="fab fa-github" />
        </a>
      </div>
    </div>
  );
}

export default Footer;
