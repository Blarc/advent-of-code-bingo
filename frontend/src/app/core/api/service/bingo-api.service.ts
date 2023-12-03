import {HttpClient, HttpHeaders} from '@angular/common/http';
import {Injectable, inject} from '@angular/core';

import {Observable} from 'rxjs';

import {CookiesService} from '../../services/cookies.service';
import {BingoBoardDto} from '../model/bingoBoardDto.model';
import {BingoCardDto} from '../model/bingoCardDto.model';
import {CreateBingoBoardDto} from '../model/createBingoBoardDto.model';
import {CreateBingoCardDto} from '../model/createBingoCardDto.model';
import {UserDto} from '../model/userDto.model';

@Injectable({
    providedIn: 'root'
})
export class BingoApiService {
    private http = inject(HttpClient);
    private cookiesService = inject(CookiesService);

    private readonly baseUrl = '/api/v1';

    private getAuthHeaders(): HttpHeaders | undefined {
        const token = this.cookiesService.getCookie('token');

        if (!token) {
            return;
        }

        return new HttpHeaders({
            'Content-Type': 'application/json',
            Authorization: `Bearer ${decodeURIComponent(token)}`
        });
    }

    public createBingoBoard(newBingoBoard: CreateBingoBoardDto): Observable<UserDto> {
        const headers = this.getAuthHeaders();
        return this.http.post<UserDto>(`${this.baseUrl}/bingoBoard`, newBingoBoard, {headers: headers});
    }

    public getBingoBoard(uuid: string): Observable<BingoBoardDto> {
        const headers = this.getAuthHeaders();
        return this.http.get<BingoBoardDto>(`${this.baseUrl}/bingoBoard/${uuid}`, {headers: headers});
    }

    public deleteBingoBoard(uuid: string): Observable<UserDto> {
        const headers = this.getAuthHeaders();
        return this.http.delete<UserDto>(`${this.baseUrl}/bingoBoard/${uuid}`, {headers: headers});
    }

    public joinBingoBoard(uuid: string): Observable<UserDto> {
        const headers = this.getAuthHeaders();
        return this.http.post<UserDto>(`${this.baseUrl}/bingoBoard/${uuid}/join`, {}, {headers: headers});
    }

    public leaveBingoBoard(uuid: string): Observable<UserDto> {
        const headers = this.getAuthHeaders();
        return this.http.delete<UserDto>(`${this.baseUrl}/bingoBoard/${uuid}/leave`, {headers: headers});
    }

    public getAllBingoCards(): Observable<BingoCardDto[]> {
        const headers = this.getAuthHeaders();
        const httpOptions = headers ? {headers} : undefined;
        return this.http.get<BingoCardDto[]>(`${this.baseUrl}/bingoCard`, httpOptions);
    }

    public createBingoCard(newBingoCard: CreateBingoCardDto): Observable<BingoCardDto> {
        return this.http.post<BingoCardDto>(`${this.baseUrl}/bingoCard`, newBingoCard);
    }

    public getUserInfo(): Observable<UserDto> {
        const headers = this.getAuthHeaders();
        return this.http.get<UserDto>(`${this.baseUrl}/me`, {headers: headers});
    }

    public clickBingoCard(uuid: string): Observable<UserDto> {
        const headers = this.getAuthHeaders();
        return this.http.post<UserDto>(`${this.baseUrl}/me/bingoCard/${uuid}`, null, {headers: headers});
    }
}
