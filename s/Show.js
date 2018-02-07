// Show component.
const Show = {
  props: ['path', 'username'],
  data() {
    return {
      key: '',
      url: '',
      owner: '',
      useCount: 0,
      found: false,
      copyHint: '复制到剪贴板',
    };
  },
  computed: {
    shortenedLink() {
      return 'http://algotest/' + this.key;
    },
  },
  methods: {
    init() {
      this.url = '';
      this.owner = '';
      this.useCount = 0;
      this.key = this.path;
      this.found = false;
      this.maybeRetrieve();
    },
    maybeRetrieve() {
      if (!!this.key) {
        this.$http
          .get('/a/manage', {
            params: {
              key: this.key,
            },
          })
          .then(r => r.json(), e => e.json())
          .then(j => this.renderResponse(j))
          .catch(e => this.myalert(e));
      }
    },
    update() {
      this.$http
        .post('/a/manage', null, {
          params: {
            key: this.key,
            url: this.url,
          },
        })
        .then(r => r.json(), e => e.json())
        .then(j => this.renderResponse(j))
        .catch(e => this.myalert(e));
    },
    remove() {
      this.$http
        .delete('/a/manage', {
          params: {
            key: this.key,
            url: this.url,
          },
        })
        .then(r => r.json(), e => e.json())
        .then(this.init)
        .catch(e => this.myalert(e));
    },
    renderResponse(j) {
      if (!j.ok) {
        this.found = false;
        this.myalert(j);
        return;
      }
      if (!j.link) {
        // not yet created, do nothing
        return;
      }
      this.found = true;
      this.url = j.link.url;
      this.owner = j.link.owner || '';
      this.useCount = j.link.useCount || 0;
      this.key = j.link.key;
      if (!this.path) {
        // created a goto link
        this.$router.push({
          name: 'show',
          props: {
            path: j.link.key,
          },
        });
      }
    },
    myalert(m) {
      console.error(JSON.stringify(m));
      if (m && m.errMsg) {
        alert(m.errMsg);
      } else {
        alert(m);
      }
    },
    onCopySuccess() {
      this.copyHint = '已复制';
    },
    onCopyError() {
      this.copyHint = '哪里不太对';
    },
  },
  created() {
    this.init()
  },
  watch: {
    path(val, prevVal) {
      if (!prevVal && !!val && this.found) {
        // just created a goto link, no need to refresh
        return;
      }
      this.init();
    },
  },
  template: `
<div>
  <div>
    <p v-if="found"><b>{{key}}</b> | 使用计数 {{useCount}} | 创建者{{owner || '不知道是谁，要不你点一下更新？'}}</p>
    <div class="form-group">
      <div v-if="found" class="input-group">
        <input v-model="shortenedLink" class="form-input" readonly />
        <button class="btn input-group-btn" 
                v-clipboard:copy="shortenedLink"
                v-clipboard:success="onCopySuccess"
                v-clipboard:error="onCopyError">
          {{copyHint}}
        </button>
      </div>
      <div v-else class="form-group">
        <input :class="['form-input', {disabled:found}]" :disabled="found" v-model="key"
               placeholder="缩写后的链接，可不填（随机分配），由大小写字母、数字及点“.”横“-”下划线“_”组成"/>
      </div>
    </div>
    <div class="form-group">
      <input class="form-input" v-model="url" placeholder="原链接"/>
    </div>
    <div class="form-group">
      <button class="btn btn-primary input-group-btn" :disabled="found&&!!owner&&username!==owner" @click="update">
        {{found?'更新':'创建'}}
      </button>
      <button v-if="found&&!!owner&&username===owner" class="btn btn-primary input-group-btn" @click="remove">
        删除
      </button>
    </div>
  </div>
</div>`,
};
