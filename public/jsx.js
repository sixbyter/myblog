var Blog = React.createClass({
  render: function() {
    var pairs = window.location.hash.substring(1).split("&"),
        hash2obj = {},
        pair,
        i;
    for ( i in pairs ) {
        if ( pairs[i] === "" ) continue;

        pair = pairs[i].split("=");
        hash2obj[ decodeURIComponent( pair[0] ) ] = decodeURIComponent( pair[1] );
    }
    if (hash2obj.article === undefined) {
      return (
        <ArticleList />
      );
    };
    console.log(hash2obj.article);
    return (
      <Article name={hash2obj.article} />
    );

  }
});

var ArticleList = React.createClass({
  componentDidMount: function() {
    $.ajax({
      url: '/api/articles',
      dataType: 'json',
      cache: false,
      type: "GET",
      success: function(data) {
        console.log(data);
        this.setState({data: data});
      }.bind(this),
      error: function(xhr, status, err) {
        console.error(status, err.toString());
      }.bind(this)
    });
  },
  getInitialState: function() {
    return {data: []};
  },
  render: function() {
    var articlelist = this.state.data.map(function (article) {
      return (
        <Title title={article.date + " " + article.title } filename={article.filename} />
      );
    });
    return (
      <div>
        {articlelist}
      </div>
    );
  }
});

var Title = React.createClass({
  handleClick: function(e) {
    window.location.href = window.location.origin + "#article=" + encodeURI(this.props.filename);
    window.React.render(
      <Article name={this.props.filename} />,
      document.getElementById('blog')
    );
  },
  render: function() {
    return(
      <div><a href="javascript:;" onClick={this.handleClick}>{this.props.title}</a></div>
    );
  }
});


var Article = React.createClass({
  componentDidMount: function() {
    var article_name = this.props.name
    $.ajax({
      url: '/api/article',
      dataType: 'json',
      cache: false,
      data: "name="+article_name,
      type: "GET",
      success: function(data) {
        console.log(data);
        this.setState({data: data});
      }.bind(this),
      error: function(xhr, status, err) {
        console.error(status, err.toString());
      }.bind(this)
    });
  },
  getInitialState: function() {
    return {data: {"content":"..."}};
  },
  rawMarkup: function() {
    var rawMarkup = marked(this.state.data.content, {sanitize: true});
    return { __html: rawMarkup };
  },
  render: function() {
    return (
      <div>
        <GoToList />
        <span dangerouslySetInnerHTML={this.rawMarkup()} />
      </div>
    );
  }
});

var GoToList = React.createClass({
  handleClick: function(e) {
    window.location.href = window.location.origin + "#";
    window.React.render(
      <Blog />,
      document.getElementById('blog')
    );
  },
  render: function() {
    return (
      <a href="javascript:;" onClick={this.handleClick}>返回列表</a>
    );
  }
});


React.render(
  <Blog />,
  document.getElementById('blog')
);