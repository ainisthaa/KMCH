import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Home from './components/pages/Home';
import Category from './components/pages/Category';
import Food from './components/pages/Food';
import Ingredient from './components/pages/Ingredient';

function App() {
  return (
    <BrowserRouter>
      {/* แถบเมนูอาจอยู่ตรงนี้ */}
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/category" element={<Category />} />
        <Route path="/food" element={<Food />} />
        <Route path="/ingredient" element={<Ingredient />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
