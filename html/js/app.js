Vue.createApp({
    data() {
        return {
            content:{},
            selected:null,
            search:""
        }
    },
    computed:{
        distributions() {
            return this.selected ? [this.selected] : Object.values(this.content).flat().filter(({distribution}) => {
                if (!this.search.trim())
                    return true
                return distribution.toLocaleLowerCase().includes(this.search.toLocaleLowerCase().trim())
            })
        },
        versions() {
            return this.selected?.versions ? this.selected.versions.filter(({url}) => {
                if (!this.search.trim())
                    return true
                return this.format("basename", url).toLocaleLowerCase().includes(this.search.toLocaleLowerCase().trim())
            }) : []
        }
    },
    methods:{
        format(type, value) {
            switch (type) {
                case "basename":
                    return value.split("/").slice(-1)[0]
                case "url":
                    return value.replace(/https?:\/\/(.*)\/?/, "$1")
            }
            return value
        }
    },
    async mounted() {
        for (const source of ["alpine", "arch", "debian", "fedora", "freebsd", "mint", "opensuse", "ubuntu"]) {
            Object.assign(this.content, {[source]:await fetch(` data/${source}.json`).then(res => res.json())})
        }
    }
}).mount("main")