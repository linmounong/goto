const Search = {
  props: ["q"],
  data() {
    return {
      q_: "",
      page: 1,
      per: 20,
      links: [],
      owner: "",
    };
  },
  methods: {
    init() {
      this.q_ = this.q || this.$route.query.q || "";
      this.page = parseInt(this.$route.query.page) || 1;
      this.per = parseInt(this.$route.query.per) || 20;
      this.owner = this.$route.query.owner || "";
      this.fetch();
    },
    fetch() {
      this.$http.get("/a/search", {
          params: {
            q: this.q_,
            page: this.page,
            per: this.per,
            owner: this.owner,
          },
        })
        .then(r => r.json(), e => e.json())
        .then(j => this.renderResponse(j));
    },
    renderResponse(j) {
      if (!j.ok) {
        this.alert(j.errMsg);
        return;
      }
      this.links = j.links || [];
    },
  },
  mounted() {
    this.init();
  },
  watch: {
    '$route' () {
      this.init();
    },
    page() {
      this.fetch();
    },
    per() {
      this.fetch();
    },
    q() {
      this.init();
    },
  },
  template: `
<div>
  <div v-for="link in links">
    <div class="divider"></div>
    <h6><router-link :to="{name:'show',params:{path:link.key}}">{{link.key}}</router-link></h6>
    <p><small>使用计数 {{link.useCount}} | 创建者 {{link.owner}} | <a class="text-ellipsis" :href="link.url">{{link.url}}</a></small></p>
  </div>
  <ul class="pagination">
    <li class="page-item">
      <a class="btn btn-link" @click="page-=1" :disabled="page<=1">上一页</a>
    </li>
    <li class="page-item">
      <!-- TODO(ynlin): this is not accurate -->
      <a class="btn btn-link" @click="page+=1" :disabled="links.length<per">下一页</a>
    </li>
  </ul>
</div>
        `,
};
