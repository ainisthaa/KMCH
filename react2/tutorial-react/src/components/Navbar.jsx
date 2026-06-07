import React from 'react';
import { Link, useLocation } from 'react-router-dom';

function Navbar() {
  const location = useLocation();
  const links = [
    { to: '/', label: 'Home' },
    { to: '/category', label: 'Category' },
    { to: '/food', label: 'Food' },
    { to: '/ingredient', label: 'Ingredient' },
  ];

  return (
    <nav className="bg-rose-200 text-white py-4 px-6 flex gap-6">
      {links.map(link => (
        <Link
          key={link.to}
          to={link.to}
          className={`font-semibold ${location.pathname === link.to ? 'text-black' : 'text-white'}`}
        >
          {link.label}
        </Link>
      ))}
    </nav>
  );
}

export default Navbar;