import React, { useEffect, useState } from "react";
import Navbar from "../Navbar";
import Pagination from "../Pagination";
import FoodModal from "../FoodModal";
import { AddButton, EditButton, DeleteButton, ViewButton } from "../Button";
import FoodViewModal from "../FoodViewModal";
import IngredientMultiSelect from "../IngredientMultiSelect";
import SortDropdown from "../SortDropdown";
import PriceRangeSelect from "../PriceRangeSelect";

function Food() {
  const [data, setData] = useState([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [isRefresh, setIsRefresh] = useState(true);

  const [page, setPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(5);

  const [showAddModal, setShowAddModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);

  const [name, setName] = useState("");
  const [price, setPrice] = useState("");
  const [categoryId, setCategoryId] = useState("");

  const [editItem, setEditItem] = useState(null);
  const [deleteId, setDeleteId] = useState(null);

  const [categoryList, setCategoryList] = useState([]);

  const [ingredientList, setIngredientList] = useState([]);
  const [selectedIngredientIds, setSelectedIngredientIds] = useState([]);

  const [foodDetail, setFoodDetail] = useState(null);
  const [showViewModal, setShowViewModal] = useState(false);

  const [searchWord, setSearchWord] = useState("");
  // const [fulldata, setfullData] = useState([])
  // const [displayData, setDisplayData] = useState(null)
  // const [searchIngredient, setSearchIngredient] = useState('')

  const [selectedIngs, setSelectedIngs] = useState([]);

  const [minPrice, setMinPrice] = useState("");
  const [maxPrice, setMaxPrice] = useState("");

  const [sortKey, setSortKey] = useState("");

  const fetchFood = () => {
    setLoading(true);

    const params = new URLSearchParams({
      page: String(page),
      limit: String(rowsPerPage),
      q: searchWord.trim(),
      sort: sortKey,
    });

    // ใส่ ingredient ที่เลือก
    selectedIngs.forEach((opt) => params.append("ingredient", opt.value));

    // ใส่ช่วงราคา (ถ้าเลือก)
    if (minPrice !== "") params.append("minPrice", minPrice);
    if (maxPrice !== "") params.append("maxPrice", maxPrice);

    fetch(`http://localhost:8890/food?${params.toString()}`)
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
  };

  useEffect(() => {
    fetchFood();
  }, [
    page,
    rowsPerPage,
    isRefresh,
    searchWord,
    selectedIngs,
    minPrice,
    maxPrice,
    sortKey,
  ]);

  // useEffect(() => {
  //     fetchFood();
  // }, [isRefresh]);

  useEffect(() => {
    fetch("http://localhost:8890/category?limit=1000")
      .then((res) => res.json())
      .then((res) => setCategoryList(res.data));
  }, []);

  useEffect(() => {
    fetch("http://localhost:8890/ingredient")
      .then((res) => res.json())
      .then((res) => setIngredientList(res.data));
  }, []);

  const handleDelete = () => {
    fetch(`http://localhost:8890/food/${deleteId}`, {
      method: "DELETE",
    }).then(() => {
      setShowDeleteModal(false);
      setDeleteId(null);
      fetchFood();
    });
  };

  const handleViewDetails = (foodId) => {
    console.log(foodId);
    fetch(`http://localhost:8890/food/${foodId}/ingredient`) // เรียก backend ที่ทำไว้
      .then((res) => res.json())
      .then((data) => {
        setFoodDetail(data);
        setShowViewModal(true);
      })
      .catch((err) => {
        console.error("โหลดข้อมูลอาหารล้มเหลว", err);
      });
  };

  const handleSubmit = () => {
    console.log(selectedIngredientIds);
    // console.log("This is edit item => ", editItem)
    const payload = {
      name,
      price: Number(price),
      category_id: Number(categoryId),
      ingredients_id: selectedIngredientIds,
    };
    console.log("Payload: ", payload);

    const url = editItem
      ? `http://localhost:8890/food/${editItem.id}`
      : `http://localhost:8890/food`;
    const method = editItem ? "PUT" : "POST";

    fetch(url, {
      method,
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    }).then((res) => {
      setShowAddModal(false);
      setShowEditModal(false);
      setEditItem(null);
      setName("");
      setPrice("");
      setCategoryId("");
      setIsRefresh(!isRefresh); // Trigger re-fetch
      fetchFood();
      setSelectedIngredientIds([]);
    });
  };

  const openAddModal = () => {
    setName("");
    setPrice("");
    setCategoryId("");
    setEditItem(null);
    setShowAddModal(true);
    setSelectedIngredientIds([]);
  };

  const openEditModal = (itemdata) => {
    console.log("item", itemdata);

    fetch(`http://localhost:8890/food/${itemdata.id}/ingredient`)
      .then((res) => res.json())
      .then((item) => {
        console.log("item in then", itemdata);
        setEditItem({ ...data, id: itemdata.id });
        setName(item.name);
        setPrice(item.price);
        setCategoryId(item.category_id);
        setSelectedIngredientIds(item.ingredients.map((i) => i.ingredient_id));
        setShowEditModal(true);
      })
      .catch((err) => {
        console.error("โหลดข้อมูลแก้อาหารล้มเหลว", err);
      });
  };

  return (
    <div className="min-h-screen bg-pink-100">
      <Navbar />
      <div className="flex flex-col items-center mt-8">
        <h2 className="text-2xl font-bold text-rose-800 mb-4">ตารางอาหาร</h2>
        {/* <pre>
          {JSON.stringify(displayData, null, 2)}
        </pre> */}
        <div className="flex gap-30">
          <SortDropdown
            sortKey={sortKey}
            setSortKey={setSortKey}
            setPage={setPage}
          />
          <div>
            <p>Search food name</p>
            <input
              className="bg-amber-200  "
              onChange={(e) => setSearchWord(e.target.value)}
            />
            {/* <div className="min-w-[260px]"> */}
          </div>
          <div>
            <p>Select ingredients</p>
            <IngredientMultiSelect
              value={selectedIngs}
              onChange={(vals) => {
                setSelectedIngs(vals);
                setPage(1);
              }}
            />
          </div>
          <PriceRangeSelect
            setMaxPrice={setMaxPrice}
            setMinPrice={setMinPrice}
            setPage={setPage}
          />
        </div>

        <div className="w-4/5 flex justify-end mb-2">
          <AddButton onClick={openAddModal} />
        </div>
        {/* <pre>
          {JSON.stringify(displayData, null, 2)}
        </pre> */}
        {loading ? (
          <p>กำลังโหลด...</p>
        ) : error ? (
          <p className="text-red-500">{error}</p>
        ) : (
          <>
            <table className="w-4/5 border border-gray-300 shadow-md text-center">
              <thead className="bg-rose-300 text-white">
                <tr>
                  <th className="py-2">No.</th>
                  <th>ชื่ออาหาร</th>
                  <th>ราคา</th>
                  <th>Category</th>
                  <th>วัตถุดิบ</th>
                  <th>อัปเดตล่าสุด</th>
                  <th className="">
                    <div className="flex justify-between">
                      <div>การจัดกา </div>
                      <div>^</div>
                    </div>
                  </th>
                </tr>
              </thead>
              <tbody>
                {Array.isArray(data) && data.length > 0 ? (
                  data.map((item, idx) => (
                    <tr key={item.id} className="even:bg-gray-100">
                      <td>{(page - 1) * rowsPerPage + idx + 1}</td>
                      <td>{item.name}</td>
                      <td>{item.price}</td>
                      <td>
                        {categoryList.find((cat) => cat.id === item.category_id)
                          ?.name || "ไม่ทราบ"}
                      </td>
                      <td>
                        {Array.isArray(item.ingredients) &&
                        item.ingredients.length > 0
                          ? item.ingredients
                              .map((ing) => ing.ingredient_name)
                              .join(", ")
                          : "ไม่มี"}
                      </td>
                      <td>{new Date(item.updated_at).toLocaleString()}</td>
                      <td>
                        <EditButton onClick={() => openEditModal(item)} />
                        <DeleteButton
                          onClick={() => {
                            setDeleteId(item.id);
                            setShowDeleteModal(true);
                          }}
                        />
                        <ViewButton
                          onClick={() => handleViewDetails(item.id)}
                        />
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan={6} className="py-4 text-gray-500">
                      ไม่พบรายการอาหาร
                    </td>
                  </tr>
                )}
              </tbody>
            </table>

            <Pagination
              page={page}
              setPage={setPage}
              rowsPerPage={rowsPerPage}
              setRowsPerPage={setRowsPerPage}
              total={total}
            />
          </>
        )}

        <FoodModal
          show={showAddModal || showEditModal || showDeleteModal}
          isEdit={!!editItem && !showDeleteModal}
          isDelete={showDeleteModal}
          onClose={() => {
            setShowAddModal(false);
            setShowEditModal(false);
            setShowDeleteModal(false);
            setEditItem(null);
            setName("");
            setPrice("");
            setCategoryId("");
            setDeleteId(null);
            setSelectedIngredientIds([]);
          }}
          onSubmit={handleSubmit}
          onDelete={handleDelete}
          name={name}
          setName={setName}
          price={price}
          setPrice={setPrice}
          categoryId={categoryId}
          setCategoryId={setCategoryId}
          categoryList={categoryList}
          ingredientList={ingredientList}
          selectedIngredientIds={selectedIngredientIds}
          setSelectedIngredientIds={setSelectedIngredientIds}
        />
        {/* <pre>
          {JSON.stringify(foodDetail, null, 2)}
        </pre> */}
        <FoodViewModal
          foodDetail={foodDetail}
          categoryList={categoryList}
          onClose={() => {
            setShowViewModal(false);
            setFoodDetail(null);
          }}
        />
      </div>
    </div>
  );
}

export default Food;
