(function global(L) {
  const DEFAULT_LATITUDE = 35.8750625;
  const DEFAULT_LONGITUDE = -78.84066989999997;
  const DEFAULT_ZOOM = 12;
  const MIN_ZOOM = 2;
  const MAX_ZOOM = 30;
  const HEAT_MODIFIER = 1000;

  const apiRoot = 'http://localhost:7071';
  let heatLayer;

  /**
   * Initializes the map.
   */
  function init() {
    const map = L.map('map').setView([DEFAULT_LATITUDE, DEFAULT_LONGITUDE], DEFAULT_ZOOM).setMinZoom(MIN_ZOOM).setMaxZoom(MAX_ZOOM);
    L.tileLayer('https://api.tiles.mapbox.com/v4/{id}/{z}/{x}/{y}.png?access_token={accessToken}', {
      maxZoom: 18,
      id: 'mapbox.streets',
      accessToken: 'pk.eyJ1IjoicGF0cmlja3ZhbGxlIiwiYSI6ImNqb2E1ajk1bDAyNWozcm5zbm1rcHZ6NTMifQ.LqzGXiSu0KkbB3SQ7mmZuA',
    }).addTo(map);

    /**
     * Fired when a move/resize/load/etc event happens in the map.
     */
    function onEvent() {
      populateHeatmap(map);
    }
    map.on('load', onEvent);
    map.on('moveend', onEvent);
    map.on('resize', onEvent);
    map.on('zoomend', onEvent);

    // Attempt to recenter the map around the browser's geo coords.
    fetchGeoCoords().then(function(coords) {
      map.setView([coords.latitude, coords.longitude], map.getZoom());
    });
  }
  init();

  /**
   * Attempts to fetch the user's geo coordinates from the browser. If not available,
   * it defaults to local geo coordinates.
   *
   * @return {Promise}
   */
  function fetchGeoCoords() {
    return new Promise(function(resolve) {
      if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(function(position) {
          resolve({'latitude': position.coords.latitude, 'longitude': position.coords.longitude});
        }, function() {
          resolve({'latitude': DEFAULT_LATITUDE, 'longitude': DEFAULT_LONGITUDE});
        });
      } else {
        resolve({'latitude': DEFAULT_LATITUDE, 'longitude': DEFAULT_LONGITUDE});
      }
    });
  }

  /**
   * Populates the heatmap layer on the supplied map.
   *
   * @param {Object} map The map to populate.
   * @return {Promise}
   */
  function populateHeatmap(map) {
    return new Promise(function(resolve, reject) {
      const bounds = map.getBounds();
      const ne = bounds.getNorthEast();
      const sw = bounds.getSouthWest();
      fetchIPv6Data(sw.lat, ne.lat, ne.lng, sw.lng).then(function(data) {
        const maxTemp = Math.log(data.maxCount + 1);
        const coords = [];
        data.results.forEach(function(address) {
          const temp = Math.log(address.count + 1);
          const intensity = HEAT_MODIFIER * (temp / maxTemp);
          coords.push([address.latitude, address.longitude, intensity]);
        });
        if (heatLayer) {
          map.removeLayer(heatLayer);
        }
        heatLayer = L.heatLayer(coords, {
          radius: calculateRadius(map),
        }).addTo(map);
      });
    });
  }

  /**
   * Fetches IPv6 data from the API according to the supplied boundaries.
   *
   * @param {number} minLatitude The lower boundary of latitude to fetch.
   * @param {number} maxLatitude The upper boundary of latitude to fetch.
   * @param {number} minLongitude The lower boundary of longitude to fetch.
   * @param {number} maxLongitude The upper boundary of latitude to fetch.
   * @return {Promise}
   */
  function fetchIPv6Data(minLatitude, maxLatitude, minLongitude, maxLongitude) {
    return new Promise(function(resolve, reject) {
      const url = new URL(apiRoot + '/v1/ipv6');
      url.searchParams.append('minLatitude', minLatitude);
      url.searchParams.append('maxLatitude', maxLatitude);
      url.searchParams.append('minLongitude', minLongitude);
      url.searchParams.append('maxLongitude', maxLongitude);

      // Make the API call.
      fetch(url, {headers: {'Accept': 'application/x-protobuf'}}).then(function(response) {
        if (!response.ok) {
          throw new Error('received status: ' + response.status + ' ' + response.statusText);
        }
        return response.arrayBuffer();
      }).then(function(buffer) {
        protobuf.load('js/proto-bundle.json', function(err, root) {
          if (err) {
            throw err;
          }
          const ListResult = root.lookupType('ipv6.ListResult');
          const message = ListResult.decode(new Uint8Array(buffer));
          resolve(message);
        });
      }).catch(reject);
    });
  }

  /**
   * Calculates and returns the radius to use in the heatmap
   * according to the zoom level of the supplied `map`.
   *
   * @param {Object} map The map.
   * @return {number}
   */
  function calculateRadius(map) {
    if (map.getZoom() < 5) {
      return 10;
    } else if (map.getZoom() < 10) {
      return 20;
    } else if (map.getZoom() < 20) {
      return 30;
    } else {
      return 40;
    }
  }
}(window.L));
