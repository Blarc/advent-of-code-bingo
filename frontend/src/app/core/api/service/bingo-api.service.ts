import {HttpClient, HttpHeaders} from '@angular/common/http';
import {Injectable, inject} from '@angular/core';

import {Observable} from 'rxjs';

import {CookiesService} from '../../services/cookies.service';
import {BingoCardDto} from '../model/bingoCardDto';
import {UserDto} from '../model/userDto';

@Injectable({
    providedIn: 'root'
})
export class BingoApiService {
    private http = inject(HttpClient);
    private cookiesService = inject(CookiesService);

    private readonly baseUrl: string;

    constructor() {
        this.baseUrl = '/api/v1';
    }

    public getAllBingoCards(): Observable<BingoCardDto[]> {
        if (this.cookiesService.getCookie('token')) {
            const headers = new HttpHeaders({
                'Content-Type': 'application/json',
                Authorization: `Bearer ${decodeURIComponent(this.cookiesService.getCookie('token'))}`
            });
            return this.http.get<BingoCardDto[]>(`${this.baseUrl}/bingoCards`, {headers: headers});
        }
        return this.http.get<BingoCardDto[]>(`${this.baseUrl}/bingoCards`);
    }

    public getUserInfo(): Observable<UserDto> {
        const headers = new HttpHeaders({
            'Content-Type': 'application/json',
            Authorization: `Bearer ${decodeURIComponent(this.cookiesService.getCookie('token'))}`
        });
        return this.http.get<UserDto>(`${this.baseUrl}/me`, {headers: headers});
    }

    public clickBingoCard(id: number): Observable<UserDto> {
        const headers = new HttpHeaders({
            'Content-Type': 'application/json',
            Authorization: `Bearer ${decodeURIComponent(this.cookiesService.getCookie('token'))}`
        });
        return this.http.post<UserDto>(`${this.baseUrl}/me/bingoCard/${id}`, null, {headers: headers});
    }
}
