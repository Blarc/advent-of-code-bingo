import {AsyncPipe, NgForOf, NgIf} from '@angular/common';
import {Component, OnInit} from '@angular/core';
import {FormBuilder, FormControl, FormGroup, ReactiveFormsModule} from '@angular/forms';
import {RouterLink, RouterLinkActive, RouterOutlet} from '@angular/router';

import {Observable, tap} from 'rxjs';

import {UserDto} from '../../core/api/model/userDto.model';
import {BingoApiService} from '../../core/api/service/bingo-api.service';
import {AuthService} from '../../core/services/auth.service';

@Component({
    standalone: true,
    selector: 'app-bingo-boards',
    styleUrls: ['bingo-boards.page.component.scss'],
    imports: [ReactiveFormsModule, NgIf, AsyncPipe, NgForOf, RouterLink, RouterOutlet, RouterLinkActive],
    templateUrl: './bingo-boards.page.component.html'
})
export class BingoBoardsPageComponent implements OnInit {
    public bingoBoardForm: FormGroup;

    private user?: UserDto;
    public user$?: Observable<UserDto | null>;

    constructor(
        private authService: AuthService,
        private formBuilder: FormBuilder,
        private bingoApiService: BingoApiService
    ) {
        this.bingoBoardForm = this.formBuilder.group({
            boardName: new FormControl('')
        });
    }

    ngOnInit() {
        this.user$ = this.authService.userSubject$.pipe(tap(u => u && (this.user = u)));
    }

    public joinBingoBoard() {
        const boardCode = this.bingoBoardForm.controls['boardName'].value;
        if (boardCode) {
            this.bingoApiService.joinBingoBoard(boardCode).subscribe({
                next: user => this.authService.updateUser(user)
            });
        }
    }

    public createBingoBoard() {
        if (this.user) {
            this.bingoApiService.createBingoBoard({name: this.user?.name}).subscribe({
                next: user => this.authService.updateUser(user)
            });
        }
    }

    public deleteBingoBoard() {
        if (this.user?.personal_bingo_board) {
            this.bingoApiService.deleteBingoBoard(this.user.personal_bingo_board.short_uuid).subscribe({
                next: user => this.authService.updateUser(user)
            });
        }
    }

    public leaveBingoBoard(uuid: string) {
        if (this.user) {
            this.bingoApiService.leaveBingoBoard(uuid).subscribe({
                next: user => this.authService.updateUser(user)
            });
        }
    }
}
