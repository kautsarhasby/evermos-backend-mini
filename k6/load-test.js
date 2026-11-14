import http from "k6/http";
import { check, sleep, group } from "k6";

export const options = {
  stages: [
    { duration: "20s", target: 10 },
    { duration: "40s", target: 10 },
    { duration: "10s", target: 0 },
  ],
  thresholds: {
    http_req_duration: ["p(95)<500"],
    http_req_failed: ["rate<0.01"],
    "http_req_duration{group:::Upload API}": ["p(95)<2000"],
  },
};
const fileBinary = open("test.jpg", "b");
const BASE_URL = "http://host.docker.internal:8000";

export function setup() {
  // User
  let loginUserRes = http.post(
    `${BASE_URL}/auth/login`,
    JSON.stringify({
      username: "hasby09", // depends on account
      password: "tes123123", // depend on account
    }),
    {
      headers: { "Content-Type": "application/json" },
    },
  );

  // Admin
  let loginAdminRes = http.post(
    `${BASE_URL}/auth/login`,
    JSON.stringify({
      username: "admin", // depends on account
      password: "admin123", // depends on account
    }),
    {
      headers: { "Content-Type": "application/json" },
    },
  );

  let userToken =
    loginUserRes.json("data.token") || loginUserRes.json("token") || "";
  let admToken =
    loginAdminRes.json("data.token") || loginAdminRes.json("token") || "";

  // REGION (agar terhindar dari limit hehehe)
  group("Public - Region API", function () {
    let res = http.get(`${BASE_URL}/provcity/listprovincies`);
    check(res, { "get provinces status 200": (r) => r.status === 200 });
  });

  return {
    userToken: userToken,
    adminToken: admToken,
  };
}

export default function (data) {
  const userHeaders = {
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${data.userToken}`,
    },
  };

  const uploadHeaders = {
    headers: {
      Authorization: `Bearer ${data.userToken}`,
    },
  };

  const adminHeaders = {
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${data.adminToken}`,
    },
  };

  // ADDRESS
  group("Address API", function () {
    let addressPayload = JSON.stringify({
      receiver_name: `Penerima_${__VU}`,
      phone: `081234567${__VU}${Math.floor(Math.random() * 100)}`,
      address_name: `Jl. Load Test No. ${__VU}, Bogor`,
      city: "Bogor",
      province: "Jawa Barat",
      subdistrict: "Ciwaringin",
      ward: "Ciwaringin Kaum",
      address_detail: "ya",
      user_id: 1,
    });

    let addressRes = http.post(
      `${BASE_URL}/api/address`,
      addressPayload,
      userHeaders,
    );

    if (addressRes.status !== 200 && addressRes.status !== 201) {
      console.log(
        `CREATE ADDRESS ERROR [${addressRes.status}]: ${addressRes.body}`,
      );
    }

    check(addressRes, {
      "create address status 201 or 200": (r) =>
        r.status === 201 || r.status === 200,
    });
  });

  // PRODUCT
  group("Products API", function () {
    // Create
    let productPayload = JSON.stringify({
      product_name: `Produk_${__VU}_${Date.now()}`,
      price: 50000,
      stock: 100,
      description: "Deskripsi produk load test",
      category_id: 1954,
    });

    let resCreate = http.post(
      `${BASE_URL}/api/products`,
      productPayload,
      userHeaders,
    );
    if (resCreate.status !== 200 && resCreate.status !== 201) {
      console.log(
        `CREATE PRODUCT ERROR [${resCreate.status}]: ${resCreate.body}`,
      );
    }
    check(resCreate, {
      "create product status 201 or 200": (r) =>
        r.status === 201 || r.status === 200,
    });

    // Get
    let resGet = http.get(`${BASE_URL}/api/products`, userHeaders);
    if (resGet.status !== 200 && resGet.status !== 201) {
      console.log(`GET PRODUCT ERROR [${resGet.status}]: ${resGet.body}`);
    }

    check(resGet, { "get products status 200": (r) => r.status === 200 });
  });

  // ENDPOINT PROFILE
  group("Protected - Profile API", function () {
    let res = http.get(`${BASE_URL}/auth/profile`, userHeaders);
    check(res, { "get profile status 200": (r) => r.status === 200 });
  });

  // ENDPOINT ADMIN
  group("Admin Only - Categories & Accounts", function () {
    let uniqueCategoryName = `Kategori_${__VU}_${Date.now()}`;

    let catRes = http.post(
      `${BASE_URL}/api/categories`,
      JSON.stringify({
        category_name: uniqueCategoryName,
      }),
      adminHeaders,
    );

    check(catRes, {
      "admin create category status 201 or 200": (r) =>
        r.status === 201 || r.status === 200,
    });

    let accRes = http.get(`${BASE_URL}/api/accounts`, adminHeaders);

    if (accRes.status !== 200) {
      console.log(`GET ACCOUNTS ERROR [${accRes.status}]: ${accRes.body}`);
    }

    check(accRes, { "admin get accounts status 200": (r) => r.status === 200 });
  });

  // UPLOAD ENDPOINT
  group("Upload API", function () {
    let formData = {
      file: http.file(fileBinary, "test.jpg", "image/jpeg"),
    };

    let res = http.post(`${BASE_URL}/api/upload`, formData, uploadHeaders);

    if (res.status !== 200 && res.status !== 201) {
      console.log(`UPLOAD ERROR [${res.status}]: ${res.body}`);
    }

    check(res, {
      "upload file status 200 or 201": (r) =>
        r.status === 200 || r.status === 201,
    });
  });

  // PAGINATION & FILTERING PRODUK
  group("Products - Pagination & Filtering API", function () {
    let res = http.get(
      `${BASE_URL}/api/products?page=1&limit=10&search=baju&category_id=1`,
      userHeaders,
    );

    check(res, {
      "get filtered paginated products status 200": (r) => r.status === 200,
    });
  });

  // TRANSACTIONS & PRODUCT LOGS
  group("Transactions API (Trigger Product Logs)", function () {
    let transactionPayload = JSON.stringify({
      address_id: 2,
      items: [
        {
          product_id: 1,
          quantity: 2,
        },
      ],
    });

    let res = http.post(
      `${BASE_URL}/api/transactions`,
      transactionPayload,
      userHeaders,
    );

    if (res.status !== 200 && res.status !== 201) {
      console.log(`TRANSACTION ERROR [${res.status}]: ${res.body}`);
    }

    check(res, {
      "create transaction status 201 or 200": (r) =>
        r.status === 201 || r.status === 200,
    });
  });

  sleep(1);
}
