import {Injectable, inject} from '@angular/core';

import {BehaviorSubject, filter, switchMap} from 'rxjs';

import {UserDto} from '../api/model/userDto';
import {BingoApiService} from '../api/service/bingo-api.service';
import {CookiesService} from './cookies.service';
import {RefreshService} from './refresh.service';

@Injectable({providedIn: 'root'})
export class AuthService {
    private cookieService = inject(CookiesService);
    private refreshService = inject(RefreshService);
    private bingoApiService = inject(BingoApiService);

    private tokenSubject = new BehaviorSubject<string | null>(null);
    private userSubject = new BehaviorSubject<UserDto | null>(null);

    constructor() {
        this.initializeTokenSubject();
    }

    private initializeTokenSubject(): void {
        const token = this.cookieService.getCookie('token');
        this.tokenSubject.next(token);
        this.tokenSubject
            .pipe(
                filter(token => token !== null),
                switchMap(() => this.bingoApiService.getUserInfo())
            )
            .subscribe(user => this.userSubject.next(user));
    }

    public get tokenSubject$() {
        return this.tokenSubject.asObservable();
    }

    public get userSubject$() {
        return this.userSubject.asObservable();
    }

    public isAuthenticated() {
        return this.tokenSubject.value !== null && this.userSubject.value !== null;
    }

    public logout() {
        this.cookieService.deleteCookie('token');
        this.tokenSubject.next(null);
        this.userSubject.next(null);
        this.refreshService.shouldRefreshBingoCards();
    }
}
