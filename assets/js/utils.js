// see https://stackoverflow.com/questions/49236100/copy-text-from-span-to-clipboard
function copyToClipboard(value) {
    const textArea = document.createElement('textarea');
    textArea.value = value;
    document.body.appendChild(textArea);
    textArea.select();
    textArea.setSelectionRange(0, 99999); // see https://www.w3schools.com/howto/howto_js_copy_clipboard.asp
    document.execCommand('copy');
    textArea.remove();
}

function convertLatLngToDMS(latitude, longitude, delimiter) {
    const lat = convertCoordinateToDMS(latitude);
    const lng = convertCoordinateToDMS(longitude);
    const degSymbol = '° ';
    return lat.degrees + degSymbol + lat.minutes + '\' ' + lat.seconds + '" ' + (latitude >= 0.0 ? 'N' : 'S')
        + delimiter + ' '
        + lng.degrees + degSymbol + lng.minutes + '\' ' + lng.seconds + '" ' + (longitude >= 0.0 ? 'E' : 'W');
}

// see https://stackoverflow.com/questions/37893131/how-to-convert-lat-long-from-decimal-degrees-to-dms-format
function convertCoordinateToDMS(coordinate) {
    const absolute = Math.abs(coordinate);
    const degrees = Math.floor(absolute);
    const minutesNotTruncated = (absolute - degrees) * 60;
    const minutes = Math.floor(minutesNotTruncated);
    const seconds = Number(((minutesNotTruncated - minutes) * 60.0).toFixed(1));
    return {degrees, minutes, seconds};
}

// see https://stackoverflow.com/questions/43167417/calculate-distance-between-two-points-in-leaflet
function calcDistance(from, to) {
    // returns distance in meters
    const degToRad = (deg) => deg * Math.PI / 180.0;
    const
        lat1 = degToRad(from[0]),
        lon1 = degToRad(from[1]),
        lat2 = degToRad(to[0]),
        lon2 = degToRad(to[1]);

    const dtLat = lat2 - lat1;
    const dtLon = lon2 - lon1;

    const a = Math.pow(Math.sin(dtLat * 0.5), 2) + Math.cos(lat1) * Math.cos(lat2) * Math.pow(Math.sin(dtLon * 0.5), 2);
    const c = 2 * Math.asin(Math.sqrt(a));
    const earthRadius = 6371.0 * 1000.0;
    return c * earthRadius;
}

// https://gist.github.com/demonixis/5188326
function toggleFullscreen(event) {
    let element = document.body;
    if (event instanceof HTMLElement) {
        element = event;
    }

    const isFullscreen = document.webkitIsFullScreen || document.mozFullScreen || false;
    element.requestFullScreen = element.requestFullScreen || element.webkitRequestFullScreen || element.mozRequestFullScreen || function () { return false; };
    document.cancelFullScreen = document.cancelFullScreen || document.webkitCancelFullScreen || document.mozCancelFullScreen || function () { return false; };
    isFullscreen ? document.cancelFullScreen() : element.requestFullScreen();
}

function getMapboxToken() {
    return 'qj/dxK0Hknh[2K0cYC4bFm5[VvhMBKiHknhX3uscVb4cVr5Ll[ieUK3ckevc{Xxelex[BK8/WS`3n6x,,TYcJ824h{6wwP';
}

function getBingMapsKey() {
    return '@qVCq{KpBd0{6qfNnfb5gpWpmKmecUjmfGGw{x`XBclctRD@P^6lyNE^5RU2q@7h';
}

function getRandomGreeting() {
    const greetings = [
        'Ahlan',
        'At Your Service',
        'Ahoj',
        'Ciao',
        'Fly, you fools',
        'Hallo',
        'Hello',
        'Hai',
        'Hei',
        'Hej',
        'Hey',
        'Hi',
        'Hoi',
        'Hola',
        'Nǐ hǎo',
        'Olá',
        'Selam',
        'Salut',
        'Salve',
        'Yassou',
        'Yā, Yō',
    ];
    return greetings[Math.floor(Math.random() * greetings.length)] + '!';
}

function getRandomWatermark() {
    const watermarks = [
        'https://media.giphy.com/media/1YLcZOlQTKRmo/giphy.gif',
        'https://media.giphy.com/media/e5kbmb3wX3J1S/giphy.gif',
        'https://media.giphy.com/media/xRJZH4Ajr973y/giphy.gif',
        '/assets/img/github-mark-120px-plus.png',
        '/assets/img/github-mark-light-120px-plus.png',
    ];
    return watermarks[Math.floor(Math.random() * watermarks.length)];
}

function getRandomAirport() {
    const airports = [
        [65.653953, -18.073572], // BIAR
        [65.275289, -14.407159], // BIEG
        [63.981874, -22.585138], // BIKF
        [50.004356, -113.635117], // CEJ4
        [51.100849, -114.374992], // CYBW
        [49.955209, -119.379941], // CYLW
        [44.236042, -78.356010], // CYPQ
        [46.790950, -71.385296], // CYQB
        [53.025336, -122.506666], // CYQZ
        [43.631210, -79.392649], // CYTZ
        [45.458611, -73.755783], // CYUL
        [49.195454, -123.172585], // CYVR
        [51.128014, -114.002861], // CYYC
        [55.93367, -129.984955], // CZST
        [50.900799, 4.491863], // EBBR
        [50.050961, 8.588685], // EDDF
        [51.289167, 6.781183], // EDDL
        [50.870651, 7.126755], // EDDK
        [48.350891, 11.777215], // EDDM
        [60.309704, 24.943562], // EFHK
        [53.360126, -2.26935], // EGCC
        [59.351121, -2.899011], // EGEP
        [59.351201, -2.951463], // EGEW
        [51.397408, -3.337534], // EGFF
        [51.154427, -0.177063], // EGKK
        [51.468273, -0.449825], // EGLL
        [55.948236, -3.357371], // EGPH
        [52.30064, 4.756833], // EHAM
        [55.038502, -8.341938], // EIDL
        [53.423454, -6.242396], // EIDW
        [55.627106, 12.644513], // EKCH
        [49.631989, 6.218031], // ELLX
        [60.193245, 11.10342], // ENGM
        [59.645973, 17.925159], // ESSA
        [-33.969231, 18.599907], // FACT
        [-26.145187, 28.229308], // FAOR
        [30.110722, 31.402348], // HECA
        [0.044560, 32.441913], // HUEN
        [39.220909, -106.866142], // KASE
        [33.642242, -84.441902], // KATL
        [46.155399, -123.880859], // KAST
        [39.858105, -104.674614], // KDEN
        [32.895115, -97.043709], // KDFW
        [59.349537, -151.927887], // KEB
        [40.640427, -73.782799], // KJFK
        [33.940392, -118.399071], // KLAX
        [28.435877, -81.312119], // KMCO
        [45.585381, -122.577904], // KPDX
        [41.976273, -87.905952], // KORD
        [43.994679, -88.555060], // KOSH
        [47.493300, -122.217610], // KRNT
        [32.733463, -117.200668], // KSAN
        [47.442154, -122.302376], // KSEA
        [34.849049, -111.791328], // KSEZ
        [37.618343, -122.382042], // KSFO
        [37.954254, -107.902718], // KTEX
        [45.421921, -123.802681], // KTMK
        [37.237545, -115.811913], // KXTA
        [41.418713, 19.715334], // LATI
        [41.286148, 2.074367], // LEBL
        [40.489525, -3.563492], // LEMD
        [39.547491, 2.736262], // LEPA
        [44.829159, -0.701476], // LFBD
        [45.396843, 6.632939], // LFLJ
        [48.999985, 2.54299], // LFPG
        [37.92857, 23.93948], // LGAV
        [47.431108, 19.259556], // LHBP
        [40.902412, 9.518719], // LIEO
        [41.790567, 12.243056], // LIRF
        [45.615388, 8.722326], // LIMC
        [50.105625, 14.263491], // LKPR
        [32.003166, 34.876381], // LLBG
        [47.258415, 11.352452], // LOWI
        [48.115978, 16.571272], // LOWW
        [32.693853, -16.776296], // LPMA
        [37.748543, -25.710541], // LPPD
        [38.770672, -9.136387], // LPPT
        [46.235291, 6.11266], // LSGG
        [36.897442, 30.796656], // LTAI
        [36.152981, -5.346678], // LXGB
        [14.061006, -87.218646], // MHTG
        [32.544254, -116.974213], // MMTJ
        [9.058223, -79.390503], // MPTO
        [8.480211, -83.589702], // MRSN
        [-37.007687, 174.786362], // NZAA
        [-43.489552, 172.537766], // NZCH
        [-46.898807, 168.104750], // NZRC
        [-45.021385, 168.740631], // NZQN
        [-41.329506, 174.810623], // NZWN
        [33.823826, 35.490906], // OLBA
        [25.245762, 55.366257], // OMDB
        [42.782394, 141.681763], // RJCC
        [34.788622, 135.442156], // RJOO
        [35.544724, 139.785782], // RJTT
        [37.445572, 126.446922], // RKSI
        [-53.781041, -67.753883], // SAWE
        [-22.815573, -43.247589], // SBGL
        [-0.127635, -78.358994], // SEQM
        [-51.687124, -57.779006], // SFAl
        [-22.816319, -43.246361], // SPGL
        [5.925572, -74.724850], // SQFQ
        [17.903753, -62.844720], // TFFJ
        [17.644494, -63.220311], // TNCS
        [32.360851, -64.701332], // TXKF
        [55.981003, 37.435276], // UUEE
        [27.687475, 86.731583], // VNLK
        [13.687966, 100.750328], // VTBS
        [27.405689, 89.421487], // VQPR
        [-6.131338, 106.657723], // WIII
        [-3.623944, 136.594537], // WX53
        [1.364033, 103.994685], // WSSS
        [-27.39019, 153.1185], // YBBN
        [-42.839382, 147.509583], // YMHB
        [-37.66573, 144.847061], // YMML
        [-31.944918, 115.968629], // YPPH
        [-33.935139, 151.176849], // YSSY
    ];
    return airports[0];
    if (airports.length > 0) {
        return airports[Math.floor(Math.random() * airports.length)];
    } else {
        return [0, 0];
    }
}

function roxy(str, key) {
    if (!key) {
        key = 1;
    }
    let output = '';
    for (let i = 0; i < str.length; ++i) {
      output += String.fromCharCode(key ^ str.charCodeAt(i));
    }
    return output;
}
