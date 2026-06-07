import React, { useEffect, useState } from "react";
import Navbar from "../Navbar";
import Pagination from "../Pagination";
import ModalComponent from "../ModalComponent";
import { AddButton, EditButton, DeleteButton } from "../Button";

function Category() {
  const [data, setData] = useState([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(5);

  const [editId, setEditId] = useState(null);
  const [editName, setEditName] = useState("");
  const [showModal, setShowModal] = useState(false);
  const [showAddModal, setShowAddModal] = useState(false);
  const [newCategoryName, setNewCategoryName] = useState("");
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [deleteId, setDeleteId] = useState(null);

  useEffect(() => {
    fetch(`http://localhost:8890/category?page=${page}&limit=${limit}`)
      .then((res) => res.json())
      .then((result) => {
        setData(result.data);
        setTotal(result.total);
        setLoading(false);
      })
      .catch((err) => {
        setError(err.message);
        setLoading(false);
      });
  }, [page, limit]);
  
  

  const openAddModal = () => {
    
    setShowAddModal(true);
  };
  const handleEdit = (item) => {
    setEditId(item.id);
    setEditName(item.name);
    setShowModal(true);
  };

  const handleSave = () => {
    fetch(`http://localhost:8890/category/${editId}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name: editName }),
    }).then(() => {
      setShowModal(false);
      setPage(1); // refresh
    });
  };

  const handleAddCategory = () => {
    fetch("http://localhost:8890/category", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name: newCategoryName }),
    })
      .then((res) => res.json())
      .then(() => {
        setShowAddModal(false);
        setNewCategoryName("");
        setPage(1);
      });
  };

  const handleDelete = () => {
    fetch(`http://localhost:8890/category/${deleteId}`, {
      method: "DELETE",
    }).then(() => {
      setShowDeleteModal(false);
      setDeleteId(null);
      setPage(1);
    });
  };

  if (loading) return <div className="text-center mt-10">Loading...</div>;
  if (error) return <div className="text-center text-red-500">Error: {error}</div>;

  return (
    <div className="min-h-screen bg-pink-100">
      <Navbar />

      <div className="flex flex-col items-center mt-8">
        <h2 className="text-2xl font-bold text-rose-800 mb-4">ตารางหมวดหมู่</h2>
        <div className="w-4/5 flex justify-end mb-2">
        <AddButton onClick={openAddModal} />

        </div>

        <table className="w-4/5 border border-gray-300 shadow-md text-center">
          <thead className="bg-rose-300 text-white">
            <tr>
              <th className="py-2">No.</th>
              <th>ชื่อหมวดหมู่</th>
              <th>อัปเดตล่าสุด</th>
              <th>การจัดการ</th>
            </tr>
          </thead>
          <tbody>
            {data.map((item, idx) => (
              <tr key={item.id} className="even:bg-gray-100">
                <td>{(page - 1) * limit + idx + 1}</td>
                <td>{item.name}</td>
                <td>{new Date(item.updated_at).toLocaleString()}</td>
                <td>
                <EditButton onClick={() => handleEdit(item)} />
                <DeleteButton
                onClick={() => {
                  setDeleteId(item.id);
                  setShowDeleteModal(true);
                }}
              />

                </td>
              </tr>
            ))}
          </tbody>
        </table>

        <Pagination
          page={page}
          setPage={setPage}
          rowsPerPage={limit}
          setRowsPerPage={setLimit}
          total={total}
        />

      <ModalComponent
        showModal={showModal}
        setShowModal={setShowModal}
        showAddModal={showAddModal}
        setShowAddModal={setShowAddModal}
        showDeleteModal={showDeleteModal}
        setShowDeleteModal={setShowDeleteModal}
        editValue={editName}
        setEditValue={setEditName}
        onEditSubmit={handleSave}
        addValue={newCategoryName}
        setAddValue={setNewCategoryName}
        onAddSubmit={handleAddCategory}
        onDelete={handleDelete}
      />

      </div>
    </div>
  );
}

export default Category;
