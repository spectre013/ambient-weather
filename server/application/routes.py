from flask import render_template, Blueprint, request
from markupsafe import escape
from sqlalchemy.sql import func
from models import Book, Author, Genre, OrderDetail
from flask_login import current_user
from db import db
application = Blueprint('application', __name__)


@application.route('/')
def index():
    books = Book.query.order_by(func.random()).limit(50)
    return render_template('application/index.html', books=books, query='')


@application.route('/author/<slug>')
def get_authors(slug):
    slug = escape(slug)
    if slug:
        author = Author.query.filter_by(Slug=slug).first()
        if not author:
            return render_template('application/index.html')
        books = Book.query.filter_by(AuthorID=author.AuthorID).all()
        return render_template('application/authors.html', author=author, books=books, search_type='book')
    return render_template('application/index.html')


@application.route('/genre/<slug>')
def get_genre(slug):
    slug = escape(slug)
    if slug:
        genre = Genre.query.filter_by(Slug=slug).first()
        if not genre:
            return render_template('application/index.html')
        books = Book.query.filter_by(GenreID=genre.GenreID).all()
        return render_template('application/genres.html', genre=genre, books=books, search_type='genre')
    return render_template('application/index.html')


@application.route('/search')
def search():
    limit = 50
    query = escape(request.args.get('query'))
    search_type = escape(request.args.get('type'))
    books = []
    if search_type == 'book' or search_type is None:
        books = Book.query.join(Author).filter(Book.Title.ilike(f'%{query}%') | Book.Description.ilike(f'%{query}%')
                                               ).limit(limit)
    elif search_type == 'author':
        print(query)
        books = Book.query.join(Author).filter(Author.FirstName.ilike(f'%{query}%') |
                                               Author.LastName.ilike(f'%{query}%')).limit(limit)
    elif search_type == 'genre':
        books = Book.query.join(Genre).filter(Genre.GenreName.ilike(f'%{query}%')).limit(limit)
    elif search_type == 'isbn':
        books = Book.query.join(Author).filter(Book.ISBN.ilike(f'%{query}%')).limit(limit)

    bookids = []
    for book in books:
        bookids.append(book.BookID)
    orderids = []
    if current_user.is_authenticated:
        orders = db.session.query(OrderDetail).filter(OrderDetail.BookID.in_(bookids)).all()
        for order in orders:
            orderids.append(order.BookID)

    return render_template('application/index.html', books=books,
                           query=query, search_type=search_type, orderids=orderids)
