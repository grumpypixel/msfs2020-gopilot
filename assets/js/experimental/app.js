// window.CategoriesDisplayOrder = (function() {
//     const displayOrder = [
//         "breakfast",
//         "lunch",
//         "dinner",
//         "sidedish",
//         "snack",
//         "sauce",
//         "dip"
//     ];
//     return {
//         displayOrder: displayOrder
//     };
// }());
class Record {
    constructor() {
        const _ = this;
        _.timestamp = Date.now();
        _.data = {};
    }

    count() {
        return Object.keys(this.data).length;
    }

    add(key, value) {
        this.data[key] = value;
    }

    get(key) {
        return this.data.hasOwnProperty(key) ? this.data[key] : undefined;
    }
}

class Recorder {
    constructor() {
        this.records = [];
    }

    clear() {
        this.records.length = 0;
    }

    addRecord(record) {
        // TODO: test record timestamp
        if (record.count() > 0) {
            this.records.push(record);
        }
    }

    count() {
        return this.records.length;
    }
}

class MinMax {
    constructor() {
        this.data = {};
    }

    set(key, value) {
        this.data[key] = {
            ref: value,
            min: value,
            max: value,
        };
        console.log("set", key, value, typeof value)
    }

    add(key, value) {
        const current = this.data[key];
        if (current !== undefined) {
            if (value < current.ref) {
                if (value < current.min) {
                    this.data[key].min = value;
                }
            } else if (value > current.ref) {
                if (value > current.max) {
                    this.data[key].max = value;
                }
            }
        }
    }

    ref(key) {
        const current = this.data[key];
        return current !== undefined ? current.ref : undefined;
    }

    minMax(key) {
        const current = this.data[key];
        return current !== undefined ? [current.min, current.max] : [undefined, undefined];
    }

    min(key) {
        const current = this.data[key];
        return current !== undefined ? current.min : undefined;
    }

    max(key) {
        const current = this.data[key];
        return current !== undefined ? current.max : undefined;
    }

    delta(key) {
        return [this.minDelta(key), this.maxDelta(key)];
    }

    minDelta(key) {
        const current = this.data[key];
        if (key === "altitude min") {
            // console.log(key, current.ref, current.min);
        }
        return current !== undefined ? Math.abs(current.ref - current.min) : undefined;
    }

    maxDelta(key) {
        const current = this.data[key];
        if (key === "altitude max") {
            // console.log(key, current.ref, current.min);
        }
        return current !== undefined ? Math.abs(current.ref - current.max) : undefined;
    }

    print(key) {
        const current = this.data[key];
        console.log(`minmax key:${key} ref:${current.ref} min:${current.min} max:${current.max}`);
    }
}

class HeadingHelper {
    constructor() {
        this.headings = [];
        this.totalAngle = 0;
        this.begin = null;
        this.end = null;
    }

    add(heading) {
        if (heading === undefined || heading === null || typeof heading !== "number") {
            return;
        }

        this.headings.push(heading);
        const count = this.headings.length;
        if (count == 1) {
            return;
        }

        const r1 = heading;
        const r2 = this.headings[count - 2];

        let d = r1 - r2;
        if (d <= -180) { // overflow when turning right (r1 > r2)
            d = r1 + 360 - r2;
        } else if (d >= 180) { // overflow when turning left (r1 < r2)
            d = r1 - (r2 + 360);
        }

        this.totalAngle += d;
    }

    absAngle() {
        return Math.abs(this.totalAngle).toFixed(0);
    }
}

const BankState = {
    none: 0,
    inside: 1,
    wasInside: 2,
};

const BankResult = {
    none: 0,
    okay: 1,
    exceed: 2,
    subceed: 3,
}

class BankHelper {
    constructor(required, delta, minAngle) {
        this.headingHelper = new HeadingHelper();
        this.required = required;
        this.delta = delta;
        this.minAngle = minAngle;
        this.state = 0;
        this.rightTurn = null;
        this.exceeded = false;
        this.subceeded = false;
        this.min = null;
        this.max = null;
    }

    add(bank, heading) {
        if (this.state === BankState.none) {
            if (Math.abs(bank) >= this.required - this.delta) {
                this.state = BankState.inside;
                this.rightTurn = bank < 0;
                this.headingHelper.add(heading);
                // this.min = bank;
                // this.max = bank;
                console.log("bank inside!");
            }
            return;
        }

        if (this.state !== BankState.inside) {
            return;
        }

        // this.min = Math.min(this.min, bank);
        // this.max = Math.max(this.max, bank);

        this.headingHelper.add(heading);
        if (this.headingHelper.absAngle() >= this.minAngle) {
            console.log("turn okay!");
        }

        if (!this.exceeded) {
            if ((this.rightTurn && (bank < -this.required - this.delta)) ||
                (!this.rightTurn && (bank > this.required + this.delta))) {

                this.state = BankState.wasInside;
                console.log("bank outside!");

                this.exceeded = true;
                console.log("bank exceeded!");
                return
            }
        }

        if (!this.subceeded) {
            if ((this.rightTurn && (bank > -this.required + this.delta)) ||
                (!this.rightTurn && (bank < this.required - this.delta))) {

                this.state = BankState.wasInside;
                console.log("bank outside!");

                const angle = this.headingHelper.absAngle();
                if (angle < this.minAngle) {
                    this.subceeded = true;
                    console.log("bank subceed!");
                }
                return;
            }
        }
    }

    isSuccess() {
        return this.state === BankState.wasInside && this.headingHelper.absAngle() >= this.minAngle && this.subceeded === false && this.exceeded === false;
    }

    didExceed() {
        return this.exceeded === true;
    }

    didSubceed() {
        return this.subceeded === true;
    }
}

window.App = (function() {
    const app = {
        name: "STEEP TURNS X",
    }
    const dataTypes = {
        float64: "float64",
    };
    const monikers = {
        altitude: "altitude",
        heading: "heading",
        airspeed: "airspeed",
        bank: "bank",
    };
    const keys = [
        monikers.altitude,
        monikers.heading,
        monikers.airspeed,
        monikers.bank,
    ];
    const constants = {
        inputNormalClass: "input is-normal",
        inputErrorClass: "input is-danger",

        tagNormalClass: "tag is-light is-medium",
        tagPassClass: "tag is-success is-medium",
        tagFailClass: "tag is-danger is-medium",
        tagInfoClass: "tag is-info is-medium",

        noValueText: "-",
        notAvailableText: "n/a",
        evaluationInProgressText: "...",
        resultPassText: "PASS",
        resultFailText: "FAIL",

        plusMinusSign: "\u00B1",

        maxDeltaAltitude: 100,
        maxDeltaAirspeed: 10,
        maxDeltaBank: 5,
        maxDeltaHeading: 10,

        requiredBank: 45,

        requiredMinBankTurnAngle: 270,

        webSocketSupport: "WebSocket" in window,
        webSocketAddress: "ws://" + window.location.hostname + ":" + window.location.port + "/ws",
    }
    const simvars = [{
            name: "INDICATED ALTITUDE", // "PLANE ALTITUDE",
            type: dataTypes.float64,
            unit: "feet",
            moniker: monikers.altitude
        },
        {
            name: "PLANE HEADING DEGREES MAGNETIC", // MAGNETIC/TRUE
            type: dataTypes.float64,
            unit: "degrees",
            moniker: monikers.heading
        },
        {
            name: "AIRSPEED INDICATED", // AIRSPEED TRUE/INDICATED
            type: dataTypes.float64,
            unit: "knots",
            moniker: monikers.airspeed
        },
        {
            name: "PLANE BANK DEGREES",
            type: dataTypes.float64,
            unit: "degrees",
            moniker: monikers.bank
        },
        // {
        //     name: "INDICATED ALTITUDE",
        //     type: dataTypes.float64,
        //     moniker: monikers.altitudeIndicated
        // },
        // {
        //     name: "PLANE HEADING DEGREES MAGNETIC",
        //     type: dataTypes.float64,
        //     moniker: monikers.headingMagnetic
        // },
        // {
        //     name: "AIRSPEED INDICATED",
        //     type: dataTypes.float64,
        //     moniker: monikers.airspeedIndicated
        // },
    ];
    return {
        app,
        constants,
        keys,
        monikers,
        simvars,
    }
}());