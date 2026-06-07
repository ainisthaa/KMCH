import React from 'react';
import Navbar from '../Navbar';

function Home() {
  return (
    <div className="min-h-screen bg-pink-100">
      <Navbar />

      <div className="flex flex-col items-center justify-center mt-20">
        <h1 className="text-4xl font-bold text-rose-800 mb-4">ยินดีต้อนรับสู่ระบบร้านอาหาร</h1>
       
      </div>
    </div>
  );
}

export default Home;
