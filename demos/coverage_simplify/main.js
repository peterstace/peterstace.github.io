var map = L.map('map').setView([-33.867805, 151.207103], 15);
L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>',
    opacity: 0.15,
}).addTo(map);

const layers = {
  'simp_nor_0': ['#f00', 'simp_nor_0.json'],
  'simp_nor_1': ['#f00', 'simp_nor_1.json'],
  'simp_nor_2': ['#f00', 'simp_nor_2.json'],
  'simp_cov_0': ['#00f', 'simp_cov_0.json'],
  'simp_cov_1': ['#00f', 'simp_cov_1.json'],
  'simp_cov_2': ['#00f', 'simp_cov_2.json'],
};

for (const [name, url] of Object.entries(layers)) {
  const xhr = new XMLHttpRequest();
  xhr.open('GET', url[1], false);
  xhr.send();
  const data = JSON.parse(xhr.responseText);
  layers[name] = L.geoJSON(data, {
    style: {
      color: url[0],
      weight: 1,
      opacity: 0.5,
    },
    onEachFeature: (feature, layer) => {
      layer.bindTooltip(
        feature.properties.nsw_loca_2,
        {
          className: 'label',
          permanent: true,
          direction: 'center',
          opacity: 1.0,
        },
      );
    },
  });
}

for (const [name, layer] of Object.entries(layers)) {
  document.getElementById(`btn_${name}`).addEventListener('click', () => {
    map.eachLayer(layer => {
      if (layer instanceof L.GeoJSON) {
        map.removeLayer(layer);
      }
    });
    layers[name].addTo(map);
  });
}

const defaultLayer = 'simp_nor_0';
layers[defaultLayer].addTo(map);
