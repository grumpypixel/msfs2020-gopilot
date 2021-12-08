var vue = new Vue({
    el: "#app",
    constants: App.constants,
    simvars: App.simvars,
    monikers: App.monikers,
    data: {
        title: "Steep Turns Maneuver",
        subtitle: "Preview Version",
        loading: true,
        recorder: new Recorder(),
        recording: false,
        recordingStarted: false,
        sampleCount: 0,
        startRecordingTime: null,
        stopRecordingTime: null,
        recordingTime: 0,

        selectedAltitude: 1500,
        selectedHeading: 0,
        selectedAirspeed: 100,
        selectedBank: App.constants.requiredBank,

        requiredAltitude: null,
        requiredHeading: null,
        requiredAirspeed: null,
        requiredBank: App.constants.requiredBank,

        absDeltaAltitude: App.constants.maxDeltaAltitude,
        absDeltaHeading: App.constants.maxDeltaHeading,
        absDeltaAirspeed: App.constants.maxDeltaAirspeed,
        absDeltaBank: App.constants.maxDeltaBank,

        minDeltaAltitude: null,
        minDeltaAirspeed: null,

        maxDeltaAltitude: null,
        maxDeltaAirspeed: null,
        minMaxDeltaBank: null,

        fmtRequiredAltitude: null,
        fmtRequiredHeading: null,
        fmtRequiredAirspeed: null,
        fmtRequiredBank: App.constants.requiredBank,

        altitudeResult: App.constants.noValueText,
        headingResult: App.constants.noValueText,
        airspeedResult: App.constants.noValueText,
        bankResult: App.constants.noValueText,

        altitudeInputClass: App.constants.inputNormalClass,
        headingInputClass: App.constants.inputNormalClass,
        airspeedInputClass: App.constants.inputNormalClass,

        altitudeMinDeltaClass: App.constants.tagNormalClass,
        airspeedMinDeltaClass: App.constants.tagNormalClass,

        altitudeMaxDeltaClass: App.constants.tagNormalClass,
        airspeedMaxDeltaClass: App.constants.tagNormalClass,

        bankMinMaxDeltaClass: App.constants.tagNormalClass,

        headingTotalClass: App.constants.tagNormalClass,

        altitudeResultClass: App.constants.tagNormalClass,
        headingResultClass: App.constants.tagNormalClass,
        airspeedResultClass: App.constants.tagNormalClass,
        bankResultClass: App.constants.tagNormalClass,

        lastRecord: null,
        socket: null,
        timer: null,
        minmax: null,

        lastData: null,
        headingHelper: null,
        bankHelper: null,
    },
    methods: {
        applyLiveValues() {
            if (this.lastData) {
                this.selectedAltitude = Math.floor(this.lastData[App.monikers.altitude]);
                this.selectedAirspeed = Math.floor(this.lastData[App.monikers.airspeed]);
                this.selectedHeading = Math.floor(this.lastData[App.monikers.heading]);
            }
        },
        startRecording() {
            if (this.recording === true) {
                return;
            }

            this.requiredAltitude = this.selectedAltitude;
            this.requiredAirspeed = this.selectedAirspeed;
            this.requiredHeading = this.selectedHeading;
            this.requiredBank = this.selectedBank;

            this.altitudeResult = App.constants.evaluationInProgressText;
            this.headingResult = App.constants.evaluationInProgressText;
            this.airspeedResult = App.constants.evaluationInProgressText;
            this.bankResult = App.constants.evaluationInProgressText;

            this.altitudeResultClass = App.constants.tagNormalClass;
            this.headingResultClass = App.constants.tagNormalClass;
            this.airspeedResultClass = App.constants.tagNormalClass;
            this.bankResultClass = App.constants.tagNormalClass;

            const now = Date.now();
            this.startRecordingTime = now;
            this.stopRecordingTime = now;

            this.recording = true;
            this.recordingStarted = true;

            this.minmax = new MinMax();
            this.minmax.set(App.monikers.altitude, parseInt(this.requiredAltitude));
            this.minmax.set(App.monikers.airspeed, parseInt(this.requiredAirspeed));
            this.minmax.set(App.monikers.heading, parseInt(this.requiredHeading));
            this.minmax.set(App.monikers.bank, parseInt(this.requiredBank));

            this.minMaxDeltaBank = BankResult.none;

            this.headingHelper = new HeadingHelper()
            this.bankHelper = new BankHelper(App.constants.requiredBank, App.constants.maxDeltaBank, App.constants.requiredMinBankTurnAngle);

            this._startTimer();
        },
        stopRecording() {
            if (this.recording === false) {
                return;
            }

            this.recording = false;
            this.stopRecordingTime = Date.now();
            this._stopTimer();

            this._evaluateResults();
        },
        resetRecording() {
            if (this.recording === true) {
                this.stopRecording();
            }

            this.recorder.clear();
            this.sampleCount = 0;
            this.startRecordingTime = null;
            this.stopRecordingTime = null;
            this.recordingTime = 0;
            this.recordingStarted = false;
            this.lastRecord = null;

            this.requiredAltitude = null;
            this.requiredHeading = null;
            this.requiredAirspeed = null;

            this.minDeltaAltitude = null;
            this.minDeltaAirspeed = null;
            // this.minDeltaBank = null;

            this.maxDeltaAltitude = null;
            this.maxDeltaAirspeed = null;
            // this.maxDeltaBank = null;
            this.minMaxDeltaBank = null;

            this.altitudeResult = App.constants.noValueText;
            this.headingResult = App.constants.noValueText;
            this.airspeedResult = App.constants.noValueText;
            this.bankResult = App.constants.noValueText;

            this.altitudeMinDeltaClass = App.constants.tagNormalClass;
            this.airspeedMinDeltaClass = App.constants.tagNormalClass;
            this.bankMinDeltaClass = App.constants.tagNormalClass;

            this.altitudeMaxDeltaClass = App.constants.tagNormalClass;
            this.airspeedMaxDeltaClass = App.constants.tagNormalClass;
            this.bankMaxDeltaClass = App.constants.tagNormalClass;

            this.headingTotalClass = App.constants.tagNormalClass;

            this.altitudeResultClass = App.constants.tagNormalClass;
            this.headingResultClass = App.constants.tagNormalClass;
            this.airspeedResultClass = App.constants.tagNormalClass;
            this.bankResultClass = App.constants.tagNormalClass;

            this.bankMinMaxDeltaClass = App.constants.tagNormalClass;

            this.minmax = null;
            this.lastData = null;

            this.headingHelper = null;
            this.bankHelper = null;

            if (this.timer) {
                this._stopTimer();
            }
        },
        _getLastRecordedValue(key) {
            return this.lastRecord ? this.lastRecord.get(key) : undefined;
        },
        _formatValue(value, defaultValue) {
            if (defaultValue === undefined) {
                defaultValue = App.constants.noValueText;
            }
            const val = typeof value === "string" ? parseInt(value, 10) : value;
            return val !== undefined && val !== null && typeof val === "number" ? val.toFixed(0) : defaultValue;
        },
        _startTimer() {
            if (this.timer) {
                this._stopTimer();
            }
            this.timer = setInterval(() => {
                this._updateRecordingTime();
            }, 100);
        },
        _stopTimer() {
            if (this.timer) {
                clearInterval(this.timer);
                this.timer = null;
            }
        },
        _updateRecordingTime() {
            const now = Date.now();
            const elapsed = this.recording ? now - this.startRecordingTime : this.stopRecordingTime - this.startRecordingTime;
            this.recordingTime = (elapsed / 1000.0).toFixed(1);
        },
        _handleSocketEvents() {
            this.socket.onopen = () => {
                console.log("connected.");
                this._registerSimVars();
            };
            this.socket.onclose = () => {
                console.log("disconnected.");
            };
            this.socket.onmessage = (e) => {
                const msg = JSON.parse(e.data);
                this._receiveMessage(msg);
            };
        },
        _registerSimVars() {
            data = App.simvars;
            meta = "nil";
            this._sendMessage("register", data, meta);
        },
        _deregisterSimVars() {
            let data = [];
            for (v of App.simvars) {
                data.push({
                    name: v.name
                });
            }
            this._sendMessage("deregister", { data }, "");
        },
        _sendMessage(name, data, meta) {
            if (data === null || data === undefined) {
                return;
            }
            meta = meta || "";
            const msg = {
                type: name,
                data,
                meta,
                debug: 0
            };
            if (App.constants.webSocketSupport === true) {
                this.socket.send(JSON.stringify(msg));
            }
        },
        _receiveMessage(msg) {
            if (msg.hasOwnProperty("type") === false || msg.hasOwnProperty("data") === false) {
                return;
            }
            if (msg.type === "simvars") {
                this._handleSimVarsMessage(msg);
            }
        },
        _handleSimVarsMessage(msg) {
            const data = msg["data"];

            this.lastData = data;

            if (this.recording) {
                this._addRecord(data);
                this._updateMinMaxValues();

                const heading = data[App.monikers.heading];
                const bank = data[App.monikers.bank];

                this.headingHelper.add(heading);
                this.bankHelper.add(bank, heading)

                this._evaluateDeltas();
            }
            //vue.$forceUpdate();
        },
        _addRecord(data) {
            const record = new Record();
            for (key of App.keys) {
                if (data.hasOwnProperty(key) === false) {
                    continue;
                }
                const value = data[key];
                record.add(key, value);
                this.minmax.add(key, value);
            }

            this.recorder.addRecord(record);
            this.sampleCount = this.recorder.count();
            this.lastRecord = record;
        },
        _updateMinMaxValues() {
            [minAltitude, maxAltitude] = this.minmax.delta(App.monikers.altitude);
            this.minDeltaAltitude = minAltitude;
            this.maxDeltaAltitude = maxAltitude;

            [minAirspeed, maxAirspeed] = this.minmax.delta(App.monikers.airspeed);
            this.minDeltaAirspeed = minAirspeed;
            this.maxDeltaAirspeed = maxAirspeed;
        },
        _evaluateResults() {
            this.altitudeResult = App.constants.notAvailableText;
            this.headingResult = App.constants.notAvailableText;
            this.airspeedResult = App.constants.notAvailableText;
            this.bankResult = App.constants.notAvailableText;

            this.altitudeResultClass = App.constants.tagInfoClass;
            this.headingResultClass = App.constants.tagInfoClass;
            this.airspeedResultClass = App.constants.tagInfoClass;
            this.bankResultClass = App.constants.tagInfoClass;

            this._evaluateAltitudeResult();
            this._evaluateAirspeedResult();
            this._evaluateHeadingResult();
            this._evaluateBankResult();
        },
        _evaluateAltitudeResult() {
            const key = App.monikers.altitude;
            const ref = this.minmax.ref(key);
            const min = this.minmax.min(key);
            const max = this.minmax.max(key);
            const abs = this.absDeltaAltitude;

            const minOk = this._evaluateDeltaValue(ref, min, abs);
            const maxOk = this._evaluateDeltaValue(ref, max, abs);

            this.altitudeResultClass = minOk && maxOk ? App.constants.tagPassClass : App.constants.tagFailClass;
            this.altitudeResult = minOk && maxOk ? App.constants.resultPassText : App.constants.resultFailText;
        },
        _evaluateAirspeedResult() {
            const key = App.monikers.airspeed;
            const ref = this.minmax.ref(key);
            const min = this.minmax.min(key);
            const max = this.minmax.max(key);
            const abs = this.absDeltaAirspeed;

            const minOk = this._evaluateDeltaValue(ref, min, abs);
            const maxOk = this._evaluateDeltaValue(ref, max, abs);

            this.airspeedResultClass = minOk && maxOk ? App.constants.tagPassClass : App.constants.tagFailClass;
            this.airspeedResult = minOk && maxOk ? App.constants.resultPassText : App.constants.resultFailText;
        },
        _evaluateHeadingResult() {
            if (this.headingHelper === null) {
                return;
            }

            const ok = Math.abs(this.headingHelper.absAngle() - 360) <= App.constants.maxDeltaHeading;
            this.headingResultClass = ok ? App.constants.tagPassClass : App.constants.tagFailClass;
            this.headingResult = ok ? App.constants.resultPassText : App.constants.resultFailText;
        },
        _evaluateBankResult() {
            if (this.bankHelper === null) {
                return;
            }

            const ok = this.bankHelper.isSuccess();
            this.bankResultClass = ok ? App.constants.tagPassClass : App.constants.tagFailClass;
            this.bankResult = ok ? App.constants.resultPassText : App.constants.resultFailText;
        },
        _evaluateDeltas() {
            this._evaluateAltitudeDelta();
            this._evaluateAirspeedDelta();
            this._evaluateHeadingDelta();
            this._evaluateBankDelta();
        },
        _evaluateAltitudeDelta() {
            const okay = App.constants.tagNormalClass;
            const fail = App.constants.tagFailClass;

            const key = App.monikers.altitude;
            const ref = this.minmax.ref(key);
            const min = this.minmax.min(key);
            const max = this.minmax.max(key);
            const abs = this.absDeltaAltitude;

            this.altitudeMinDeltaClass = this._evaluateDeltaValue(ref, min, abs) ? okay : fail;
            this.altitudeMaxDeltaClass = this._evaluateDeltaValue(ref, max, abs) ? okay : fail;
        },
        _evaluateAirspeedDelta() {
            const okay = App.constants.tagNormalClass;
            const fail = App.constants.tagFailClass;

            const key = App.monikers.airspeed;
            const ref = this.minmax.ref(key);
            const min = this.minmax.min(key);
            const max = this.minmax.max(key);
            const abs = this.absDeltaAirspeed;

            this.airspeedMinDeltaClass = this._evaluateDeltaValue(ref, min, abs) ? okay : fail;
            this.airspeedMaxDeltaClass = this._evaluateDeltaValue(ref, max, abs) ? okay : fail;
        },
        _evaluateHeadingDelta() {
            if (this.headingHelper === null) {
                return;
            }
            const angle = this.headingHelper.absAngle();
            if (angle > 360 + App.constants.maxDeltaHeading) {
                this.headingTotalClass = App.constants.tagFailClass;
            }
        },
        _evaluateBankDelta() {
            if (this.bankHelper === null) {
                return;
            }

            if (this.bankHelper.isSuccess()) {
                this.minMaxDeltaBank = BankResult.okay;
                this.bankMinMaxDeltaClass = App.constants.tagNormalClass;
            } else if (this.bankHelper.state === BankState.inside) {
                this.minMaxDeltaBank = BankResult.okay;
                this.bankMinMaxDeltaClass = App.constants.tagNormalClass;
            } else if (this.bankHelper.didExceed()) {
                this.minMaxDeltaBank = BankResult.exceed;
                this.bankMinMaxDeltaClass = App.constants.tagFailClass;
            } else if (this.bankHelper.didSubceed()) {
                this.minMaxDeltaBank = BankResult.subceed;
                this.bankMinMaxDeltaClass = App.constants.tagFailClass;
            }
        },
        _evaluateDeltaValue(ref, val, abs) {
            return Math.abs(ref - val) <= abs;
        },
    },
    watch: {
        selectedAltitude(val) {
            this.altitudeInputClass = !isNaN(parseInt(val)) ? App.constants.inputNormalClass : App.constants.inputErrorClass;
            this.fmtRequiredAltitude = this._formatValue(this.selectedAltitude);
        },
        selectedAirspeed(val) {
            this.airspeedInputClass = !isNaN(parseInt(val)) ? App.constants.inputNormalClass : App.constants.inputErrorClass;
            this.fmtRequiredAirspeed = this._formatValue(this.selectedAirspeed);
        },
        selectedHeading(val) {
            this.headingInputClass = !isNaN(parseInt(val)) ? App.constants.inputNormalClass : App.constants.inputErrorClass;
            this.fmtRequiredHeading = this._formatValue(this.selectedHeading);
        },
    },
    computed: {
        hasSamples() {
            return this.sampleCount > 0;
        },
        fmtAbsDeltaAltitude() {
            return `${App.constants.plusMinusSign} ${this.absDeltaAltitude}`;
        },
        fmtAbsDeltaHeading() {
            return `${App.constants.plusMinusSign} ${this.absDeltaHeading}`;
        },
        fmtAbsDeltaAirspeed() {
            return `${App.constants.plusMinusSign} ${this.absDeltaAirspeed}`;
        },
        fmtAbsDeltaBank() {
            return `${App.constants.plusMinusSign} ${this.absDeltaBank}`;
        },
        fmtCurrentAltitude() {
            return this._formatValue(this._getLastRecordedValue(App.monikers.altitude));
        },
        fmtCurrentHeading() {
            return this._formatValue(this._getLastRecordedValue(App.monikers.heading));
        },
        fmtCurrentAirspeed() {
            return this._formatValue(this._getLastRecordedValue(App.monikers.airspeed));
        },
        fmtCurrentBank() {
            let val = this._getLastRecordedValue(App.monikers.bank);
            if (val !== undefined) {
                val = Math.abs(val);
            }
            return this._formatValue(val);
        },
        fmtMinDeltaAltitude() {
            return this.minDeltaAltitude !== null ? `-${this._formatValue(this.minDeltaAltitude)}` : App.constants.noValueText;
        },
        fmtMinDeltaAirspeed() {
            return this.minDeltaAirspeed !== null ? `-${this._formatValue(this.minDeltaAirspeed)}` : App.constants.noValueText;
        },
        // fmtMinDeltaBank() {
        //     return this.minDeltaBank !== null ? `-${this._formatValue(this.minDeltaBank)}` : App.constants.noValueText;
        // },
        fmtMaxDeltaAltitude() {
            return this.maxDeltaAltitude !== null ? `+${this._formatValue(this.maxDeltaAltitude)}` : App.constants.noValueText;
        },
        fmtTotalHeading() {
            if (this.headingHelper === null) {
                return App.constants.noValueText;
            }
            const angle = this.headingHelper.absAngle();
            if (angle < 360) {
                return `${angle} of 360`;
            } else {
                return `+${angle - 360} past 360`;
            }
        },
        fmtMaxDeltaAirspeed() {
            return this.maxDeltaAirspeed !== null ? `+${this._formatValue(this.maxDeltaAirspeed)}` : App.constants.noValueText;
        },
        // fmtMaxDeltaBank() {
        //     return this.maxDeltaBank !== null ? `+${this._formatValue(this.maxDeltaBank)}` : App.constants.noValueText;
        // },
        fmtMinMaxDeltaBank() {
            if (this.minMaxDeltaBank === null) {
                return "-";
            }
            switch (this.minMaxDeltaBank) {
                case BankResult.none:
                    return "...";
                case BankResult.okay:
                    return "Good!";
                case BankResult.subceed:
                    return "Subceeded!";
                case BankResult.exceed:
                    return "Exceeded!";
            }
            return "-";
        },
        isResettable() {
            return this.recording || this.sampleCount > 0 || this.recordingTime > 0;
        },
        validInputValues() {
            return !isNaN(parseInt(this.selectedAltitude)) &&
                !isNaN(parseInt(this.selectedAirspeed)) &&
                !isNaN(parseInt(this.selectedHeading));
        },
    },
    beforeCreate() {},
    created() {
        console.log("# created");
        if (App.constants.webSocketSupport === true) {
            this.socket = new WebSocket(App.constants.webSocketAddress);
            this._handleSocketEvents();
        }
    },
    mounted() {
        console.log("# mounted");
        this.loading = false;

        this.fmtRequiredAltitude = this.selectedAltitude;
        this.fmtRequiredAirspeed = this.selectedAirspeed;
        this.fmtRequiredHeading = this.selectedHeading;

        // console.log(this.$root.$options.constants);
        // console.log(this.$root.$options.simvars);
        // console.log(this.$root.$options.monikers);
    },
    updated() {},
    destroyed() {
        if (App.constants.webSocketSupport === true) {
            this._deregisterSimVars();
            this.socket.onclose = () => {};
            this.socket.close();
        }
    },
});