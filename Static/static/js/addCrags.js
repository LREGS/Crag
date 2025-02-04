const form = document.querySelector("#addCrag_form");

async function sendData() {
  const formData = new FormData(form);
  const data = Object.fromEntries(formData);

  // Convert latitude and longitude to numbers
  data.latitude = parseFloat(data.latitude);
  data.longitude = parseFloat(data.longitude);
  // name stays as string, no conversion needed

  try {
    const response = await fetch("http://localhost:6969/api/v1/crags", {
      method: "POST",
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data),
    });
    console.log(await response.json());
  } catch (e) {
    console.error(e);
  }
}

form.addEventListener("submit", (event) => {
  event.preventDefault();
  sendData();
});