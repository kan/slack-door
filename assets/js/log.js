var Message = function(data) {
  data = data || {}
  this.user = m.prop(data.user)
  this.text = m.prop(data.text)
  this.ts = m.prop(data.ts)
};
Message.list = function() {
  var apiURL = "/log/" + m.route.param("channel") + "/" 
                       + m.route.param("year") + "/"
                       + m.route.param("month") + "/"
                       + m.route.param("day");
  console.log(apiURL);
  return m.request({ method: "GET", url: apiURL }).then(function(result){
    return result.messages;
  });
}

var Log = {
  controller: function() {
    this.messages = Message.list();
  },
  view: function(ctrl) {
    return m("div", [
        ctrl.messages().map(function(msg, index) {
          return m("p", { key: msg.ts }, msg.text);
        })
    ]);
  }
};

