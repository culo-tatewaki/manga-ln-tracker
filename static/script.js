document.addEventListener("DOMContentLoaded", async () => {
  const seriesListTag = document.getElementById("series-list");

  const addFormTag            = document.getElementById("add-form");
  const addFormIdTag          = document.getElementById("add-form-id");
  const addFormTypeTag        = document.getElementById("add-form-type");
  const addFormTitleTag       = document.getElementById("add-form-title");
  const addFormChaptersTag    = document.getElementById("add-form-chapters");
  const addFormVolumesTag     = document.getElementById("add-form-volumes");
  const addFormReleaseDateTag = document.getElementById("add-form-release-date");
  const addFormStatusTag      = document.getElementById("add-form-status");
  const addFormAuthorTag      = document.getElementById("add-form-author");
  const addFormImageTag       = document.getElementById("add-form-image");
  const addFormRatingTag      = document.getElementById("add-form-rating");

  const searchFormTag        = document.getElementById("search-form")
  const searchTitleTag       = document.getElementById("search-title")
  const searchTypeTag        = document.getElementById("search-type")
  const searchStatusTag      = document.getElementById("search-status")
  const searchRatingTag      = document.getElementById("search-rating")
  const searchReleaseDateTag = document.getElementById("search-release-date")  
  
  const res = await fetch("/getall");
  const seriesList = await res.json();
  refreshSeriesList(seriesListTag, seriesList);
  
  searchFormTag.addEventListener("submit", async (e) => {
    e.preventDefault();

    const res = await fetch("/search", {
      method: "POST",
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        type: searchTypeTag.value,
        title: searchTitleTag.value,
        track: {
          status: searchStatusTag.value,
        },
        releaseDate: parseInt(searchReleaseDateTag.value),
        rating: searchRatingTag.value,
      }),
    });

    const seriesList = await res.json();
    refreshSeriesList(seriesListTag, seriesList)
  });

  addFormTag.addEventListener("submit", async (e) => {
    e.preventDefault();

    const id = parseInt(addFormIdTag.value)
    
    const url    = id != -1 ? "/update" : "/add";
    const method = id != -1 ? "PUT" : "POST";

    const res = await fetch(url, {
      method: method,
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        id: id,
        type: addFormTypeTag.value,
        title: addFormTitleTag.value,
        track: {
          chapters: parseInt(addFormChaptersTag.value),
          volumes: parseInt(addFormVolumesTag.value),
          status: addFormStatusTag.value,
          lastUpdate: new Date(),
        },
        author: addFormAuthorTag.value,
        releaseDate: parseInt(addFormReleaseDateTag.value),
        image: addFormImageTag.value,
        rating: addFormRatingTag.value,
      }),
    });

    const series = await res.json();

    if (url == "/add") {
      seriesListTag.innerHTML += generateSeriesHtml(series);
    } else {
      const toBeUpdatedTag = seriesListTag.querySelector(`[data-id="${id}"]`);
      const tempTag = document.createElement("div");
      tempTag.innerHTML = generateSeriesHtml(series);
      const updatedTag = tempTag.firstElementChild;
      seriesListTag.replaceChild(updatedTag, toBeUpdatedTag);
    }
  });

  document.getElementById("delete-btn").addEventListener("click", async () => {
    const id = addFormIdTag.value;
    await fetch(`/delete?id=${id}`, {
      method: 'DELETE',
    });

    seriesListTag.querySelector(`[data-id="${id}"]`).remove();
  });

  document.getElementById("add-btn").addEventListener("click", () => { resetForm(addFormTag) });
  document.getElementById("clear-btn").addEventListener("click", () => { resetForm(searchFormTag) });
});

function resetForm(formTag) {
  formTag.reset();
}

function generateSeriesHtml(series) {
  return `
  <div type="button" class="card mb-3 col m-2 p-0 series-id edit-btn" style="max-width: 540px;" data-id=${series.id} onclick="addValuesToEdit(event)" data-bs-toggle="modal" data-bs-target="#form-modal">
    <div class="row g-0">
      <div class="col-md-4">
        <img src=${series.image} class="rounded-start img-fluid series-image" alt="${series.title}">
      </div>
      <div class="col-md-8">
        <div class="card-body">
          <h5 class="card-title series-title">${series.title}</h5>
          <div class="card-info mb-2">
            <p>Chapters: <span class="series-chapters">${series.track.chapters}</span></p>
            <p>Volumes: &nbsp;<span class="series-volumes">${series.track.volumes}</span></p>
          </div>
          <div class="card-info mb-2">
            <p class="card-text">Type: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<span class="series-type">${series.type}</span></p>
            <p class="card-text">Author: &nbsp;&nbsp;&nbsp;<span class="series-author">${series.author}</span></p>
            <p class="card-text">Status: &nbsp;&nbsp;&nbsp;&nbsp;<span class="series-status">${series.track.status}</span></p>
            <p class="card-text">Rating: &nbsp;&nbsp;&nbsp;<span class="series-rating">${series.rating}</span></p>
            <p class="card-text">Release: &nbsp;<span class="series-release-date">${series.releaseDate}</span></p>
          </div>
          <p class="card-text"><small class="text-muted">Last updated ${(new Date(series.track.lastUpdate)).toLocaleString()}</small></p>
        </div>
      </div>
    </div>
  </div>
  `;
}

function refreshSeriesList(seriesListTag, seriesList) {
  seriesListTag.innerHTML = "";

  seriesList.forEach((series) => {
    seriesListTag.innerHTML += generateSeriesHtml(series)
  });
}

function addValuesToEdit({ currentTarget }) {  
  document.getElementById("add-form-id").value = currentTarget.getAttribute("data-id");
  document.getElementById("add-form-type").value = currentTarget.querySelector(".series-type").textContent;
  document.getElementById("add-form-title").value = currentTarget.querySelector(".series-title").textContent;
  document.getElementById("add-form-chapters").value = currentTarget.querySelector(".series-chapters").textContent;
  document.getElementById("add-form-volumes").value = currentTarget.querySelector(".series-volumes").textContent;
  document.getElementById("add-form-release-date").value = currentTarget.querySelector(".series-release-date").textContent;
  document.getElementById("add-form-status").value = currentTarget.querySelector(".series-status").textContent;
  document.getElementById("add-form-author").value = currentTarget.querySelector(".series-author").textContent;
  document.getElementById("add-form-image").value = currentTarget.querySelector(".series-image").src;
  document.getElementById("add-form-rating").value = currentTarget.querySelector(".series-rating").textContent;
}