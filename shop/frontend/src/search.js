function getParameterByName(name, url) {
  if (!url) url = window.location.href;
  name = name.replace(/[\[\]]/g, '\\$&');
  var regex = new RegExp('[?&]' + name + '(=([^&#]*)|&|#|$)'),
      results = regex.exec(url);
  if (!results) return null;
  if (!results[2]) return '';
  return decodeURIComponent(results[2].replace(/\+/g, ' '));
}

function initialSearch(event) {
  const q = getParameterByName("q") || "*";
  const page = getParameterByName("page") || 0;
  doSearch(true, q, page);
}

function handleEnter(event) {
  if (event.code == "Enter") {
    doSearch(true);
  }
}

function doSearch(reset, inSearchTerm, inPage) {
  let searchTerm = inSearchTerm ? inSearchTerm : document.querySelector(".searchTerm").value;
  const page = inPage ? inPage : document.querySelector(".paginator select").value;
  history.pushState({
    searchTerm,
    page
  }, "Suche - "+searchTerm, "index.html?q="+searchTerm+"&page="+page);
  if (searchTerm === ""){
    searchTerm = "*";
  }
  toggleSearchForm(false);
  fetch("/api/products/search?q=" + searchTerm + "&page=" + page)
    .then(result => {
      if (result.ok) {
        return result.json();
      }
      throw new Error("Unexpected Status Code: " + result.status);
    })
    .then(result => {
      toggleSearchForm(true);
      focusInputField();
      let resultWrapper = document.querySelector(".resultWrapper");
      while (resultWrapper.firstChild) {
        resultWrapper.removeChild(resultWrapper.firstChild);
      }
      refreshPaginator(result.total, reset, page)
      if (result.hits.length > 0) {
        const resultList = generateResultList(result.hits);
        document.querySelector(".resultWrapper").appendChild(resultList);
      } else {
        document.querySelector(".resultWrapper").innerHTML =
          `<div class="null-hits">
            No results for ${searchTerm}<br />
            Please check your spelling or use a more general term.
          </div>`;
      }
    })
    .catch(err => {
      toggleSearchForm(true);
      document.querySelector(".resultWrapper").innerHTML =
        "Error during Search: <br />'" + err + "'";
      console.error(err);
    });
}

function pad(n, width, z) {
  z = z || "0";
  n = n + "";
  return n.length >= width ? n : new Array(width - n.length + 1).join(z) + n;
}

function refreshPaginator(total,reset, current){
  const paginator = document.querySelector(".paginator")
  if (total <= 10) {
    paginator.querySelector("select").value = 0
    paginator.style = "display: none;"
  }else{
    paginator.style = "display: block;"
    if(reset){
      const selector = paginator.querySelector("select")
      while (selector.firstChild) {
        selector.removeChild(selector.firstChild);
      }
      const page = parseInt(current);
      for (let i=0;i<(total/10);i++){
        const option = document.createElement("option");
        option.selected = (i === page);
        option.value = i;
        option.innerText = (i + 1);
        selector.appendChild(option)
      }
    }
  }
}

function generateResultList(results) {
  const resultList = document.createElement("div");
  resultList.classList.add("product-list");
  for (const result of results) {
    const z = Math.random();
    let available = "";
    if (z<0.3){
      if(z<0.1){
        available = "Zurzeit nicht verfügbar";
      }else{
        available = "Bregrenzter Vorrat";
      }
    }
    resultList.innerHTML += `
            <div class="product" data-product-id="${result.id}">
                <figure>
                    <img 
                        src="${result.imageUrl}" 
                        alt="Product Image" class="product-image" 
                        data-discounted="${result.discounted}" />
                </figure>
                <div class="product-description">

                    <div class="info">
                    <h1>${result.title}</h1>
                    <span class="brand">
                    ${result.brand}
                    </span>
                    <br />
                    <span class="availability">
                    ${available}
                    </span>
                    </div>
                    <div class="spacer"></div>
                    <div class="price">
                    ${Math.trunc(result.price / 100)},${pad(
      result.price % 100,
      2
    )}€
                    </div>
                </div>
            </div>
        `;
  }
  return resultList;
}

function focusInputField() {
  const inputField = document.querySelector(".searchTerm");
  inputField.focus();
  const len = inputField.value.length;
  inputField.setSelectionRange(len, len);
}

function toggleSearchForm(active) {
  if (active) {
    document.querySelector(".searchTerm").disabled = undefined;
    document.querySelector(".spinnerWrapper").style.display = "none";
  } else {
    document.querySelector(".searchTerm").disabled = true;
    document.querySelector(".spinnerWrapper").style.display = "flex";
    document.querySelector(".resultWrapper").innerHTML = "";
  }
}
