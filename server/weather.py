from application import create_app

if __name__ == "__main__":
    create_app('development').run(host='0.0.0.0', port=3000, debug=True, threaded=True)
