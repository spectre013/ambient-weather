import moment from "moment";


function Footer() {
    return (
        <footer className='copyright'>
            &copy; 2018-{moment().format('YYYY')} : zoms.net
        </footer>
    );
}

export default Footer;