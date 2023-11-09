import {HttpClient, HttpHeaders} from '@angular/common/http';
import {Injectable} from '@angular/core';

import {Observable} from 'rxjs';

import {CookiesService} from '../../services/cookies.service';
import {BingoCardDto} from '../model/bingoCardDto';
import {UserDto} from '../model/userDto';

@Injectable({
    providedIn: 'root'
})
export class ApiService {
    private readonly baseUrl: string;

    constructor(
        private http: HttpClient,
        private cookiesService: CookiesService
    ) {
        this.baseUrl = 'http://localhost:8080/api/v1';
    }

    public getAllBingoCards(): Observable<BingoCardDto[]> {
        return this.http.get<BingoCardDto[]>(`${this.baseUrl}/bingoCards`);
    }

    public clickBingoCard(id: number): Observable<UserDto> {
        const headers = new HttpHeaders({
            'Content-Type': 'application/json',
            Authorization: `Bearer ${decodeURIComponent(this.cookiesService.getCookie('token'))}`
        });
        return this.http.post<UserDto>(`${this.baseUrl}/me/bingoCard/${id}`, null, {headers: headers});
    }
}
