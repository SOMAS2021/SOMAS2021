import React from 'react';
import '../style/Button.css';
import { Link } from 'react-router-dom';

const STYLES = ['btn--primary', 'btn--outline'];

const SIZES = ['btn--medium', 'btn--large'];

export const Button = ({ children, type, onClick, buttonStyle, buttonSize }) => {
  const checkButtonStyle = STYLES.includes(buttonStyle) ? buttonStyle : STYLES[0];

  const checkButtonSize = SIZES.includes(buttonSize) ? buttonSize : SIZES[0];

  return (
    <Link to="/about" className="btn-mobile">
      <button
        className={`btn ${checkButtonSize} ${checkButtonStyle}`}
        onClick={onClick}
        type={type}
      >
        {children}
      </button>
    </Link>
  );
};
