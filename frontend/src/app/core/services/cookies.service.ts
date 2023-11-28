import {Injectable} from '@angular/core';

@Injectable({providedIn: 'root'})
export class CookiesService {
    public getCookie(name: string) {
        const cookieName = `${name}=`;
        let cookieValue = '';

        document.cookie.split(';').forEach(cookie => {
            const trimmedCookie = cookie.trim();
            if (trimmedCookie.indexOf(cookieName) === 0) {
                cookieValue = trimmedCookie.substring(cookieName.length, trimmedCookie.length);
            }
        });

        return cookieValue;
    }

    public deleteCookie(name: string) {
        this.setCookie(name, '', -1);
    }

    public setCookie(name: string, value: string, expireDays: number, path = '') {
        const expirationDate: Date = new Date();
        expirationDate.setTime(expirationDate.getTime() + expireDays * 24 * 60 * 60 * 1000);
        const expires = `expires=${expirationDate.toUTCString()}`;
        const cookiePath: string = path.length > 0 ? `; path=${path}` : '';
        document.cookie = `${name}=${value}; ${expires}${cookiePath}`;
    }
}
