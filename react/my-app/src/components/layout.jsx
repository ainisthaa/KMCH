import { Outlet } from "react-router-dom";

export const Topbar = () => {
  return (
    <div className="bg-blue-500 text-white p-4">
      <h1 className="text-3xl font-bold">Topbar</h1>
    </div>
  );
};

export const Footer = () => {
  return (
    <div className="bg-gray-800 text-white p-4">
      <h1 className="text-3xl font-bold">Footer</h1>
    </div>
  );
};

export default function Layout() {
  return (
    <div>
      <Topbar />
      <Outlet />
      <Footer />
    </div>
  );
}
